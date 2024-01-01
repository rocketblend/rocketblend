package rocketpack

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/go-git/go-git/v5"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
)

type (
	Service interface {
		Get(ctx context.Context, forceUpdate bool, references ...reference.Reference) (map[reference.Reference]*RocketPack, error)
		Remove(ctx context.Context, references ...reference.Reference) error
		Insert(ctx context.Context, packs map[reference.Reference]*RocketPack) error
	}

	Options struct {
		Logger      logger.Logger
		StoragePath string
	}

	Option func(*Options)

	getResult struct {
		packs map[reference.Reference]*RocketPack
		error error
	}

	service struct {
		logger      logger.Logger
		storagePath string
	}
)

func WithStoragePath(storagePath string) Option {
	return func(o *Options) {
		o.StoragePath = storagePath
	}
}

func WithLogger(logger logger.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}

func NewService(opts ...Option) (Service, error) {
	options := &Options{
		Logger: logger.NoOp(),
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.StoragePath == "" {
		return nil, fmt.Errorf("storage path is required")
	}

	err := os.MkdirAll(options.StoragePath, 0755)
	if err != nil {
		return nil, err
	}

	options.Logger.Debug("Initializing rocketpack service", map[string]interface{}{
		"storagePath": options.StoragePath,
	})

	return &service{
		logger:      options.Logger,
		storagePath: options.StoragePath,
	}, nil
}

func (s *service) Get(ctx context.Context, forceUpdate bool, references ...reference.Reference) (map[reference.Reference]*RocketPack, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make(chan getResult)
	var wg sync.WaitGroup
	wg.Add(len(references))

	for _, ref := range references {
		go func(ref reference.Reference) {
			defer wg.Done()

			packs, err := s.get(ctx, forceUpdate, ref)
			if err != nil {
				cancel()
				results <- getResult{packs: nil, error: err}
				return
			}

			results <- getResult{packs: packs, error: nil}
		}(ref)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	packages := make(map[reference.Reference]*RocketPack)
	for res := range results {
		if res.error != nil {
			return nil, res.error
		}

		if res.packs != nil {
			for ref, pack := range res.packs {
				packages[ref] = pack
			}
		}
	}

	return packages, nil
}

func (s *service) Remove(ctx context.Context, references ...reference.Reference) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := make(chan error, len(references))
	var wg sync.WaitGroup
	wg.Add(len(references))

	for _, ref := range references {
		go func(ref reference.Reference) {
			defer wg.Done()

			err := s.removePackage(ctx, ref)
			if err != nil {
				cancel()
				errs <- err
				return
			}
		}(ref)
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) Insert(ctx context.Context, packs map[reference.Reference]*RocketPack) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := make(chan error, len(packs))
	var wg sync.WaitGroup
	wg.Add(len(packs))

	for ref, pack := range packs {
		go func(ref reference.Reference, pack *RocketPack) {
			defer wg.Done()

			err := s.insertPackage(ctx, ref, pack)
			if err != nil {
				cancel()
				errs <- err
				return
			}
		}(ref, pack)
	}

	go func() {
		wg.Wait()
		close(errs)
	}()

	var firstErr error
	for err := range errs {
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
		}
	}

	return firstErr
}

func (s *service) insertPackage(ctx context.Context, ref reference.Reference, pack *RocketPack) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	packagePath := filepath.Join(s.storagePath, ref.String(), FileName)

	err := os.MkdirAll(filepath.Dir(packagePath), 0755)
	if err != nil {
		s.logger.Error("Error creating directory", map[string]interface{}{"error": err, "reference": ref.String(), "path": filepath.Dir(packagePath)})
		return err
	}

	if err := Save(packagePath, pack); err != nil {
		s.logger.Error("Error saving package", map[string]interface{}{"error": err, "reference": ref.String()})
		return err
	}

	return nil
}

func (s *service) get(ctx context.Context, forceUpdate bool, ref reference.Reference) (map[reference.Reference]*RocketPack, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.logger.Info("Processing reference", map[string]interface{}{"reference": ref.String()})

	packages := make(map[reference.Reference]*RocketPack)

	repo, err := ref.GetRepo()
	if err != nil {
		s.logger.Error("Error getting repository", map[string]interface{}{"error": err, "reference": ref.String()})
		return nil, err
	}

	repoPath := filepath.Join(s.storagePath, repo)
	packagePath := filepath.Join(s.storagePath, ref.String(), FileName)

	// The repository does not exist locally, clone it
	if _, err := os.Stat(repoPath); os.IsNotExist(err) || forceUpdate && !ref.IsLocalOnly() {
		repoURL, err := ref.GetRepoURL()
		if err != nil {
			s.logger.Error("Error getting repository URL", map[string]interface{}{"error": err, "reference": ref.String()})
			return nil, err
		}

		s.logger.Info("Cloning repository", map[string]interface{}{"repoURL": repoURL, "path": repoPath, "reference": ref.String()})
		if err := s.cloneRepo(ctx, repoPath, repoURL); err != nil {
			return nil, err
		}
	}

	// Check if the file exists in the repository
	if _, err := os.Stat(packagePath); os.IsNotExist(err) || forceUpdate && !ref.IsLocalOnly() {
		// The file does not exist or forced update, pull the latest changes
		s.logger.Info("Pulling latest changes for repository", map[string]interface{}{"path": packagePath, "reference": ref.String()})
		if err := s.pullChanges(ctx, repoPath); err != nil {
			return nil, err
		}
	}

	pack, err := Load(packagePath)
	if err != nil {
		s.logger.Error("Error loading package", map[string]interface{}{"error": err, "reference": ref.String(), "path": packagePath})
		return nil, err
	}

	deps := pack.GetDependencies()
	if len(deps) > 0 {
		s.logger.Debug("Package has dependencies", map[string]interface{}{"reference": ref.String()})

		// Get the dependencies
		depPackages, err := s.Get(ctx, forceUpdate, deps...)
		if err != nil {
			s.logger.Error("Error getting dependency packages", map[string]interface{}{"error": err, "reference": ref.String()})
			return nil, err
		}

		// Add the dependencies to the packages map
		for index, dep := range depPackages {
			packages[index] = dep
		}

		s.logger.Debug("Dependency packages successfully loaded", map[string]interface{}{"reference": ref.String()})
	}

	packages[ref] = pack

	return packages, nil
}

func (s *service) removePackage(ctx context.Context, reference reference.Reference) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.logger.Info("Processing reference", map[string]interface{}{"reference": reference.String()})

	repo, err := reference.GetRepo()
	if err != nil {
		s.logger.Error("Error getting repository path", map[string]interface{}{"error": err, "reference": reference.String()})
		return err
	}

	repoPath := filepath.Join(s.storagePath, repo)

	// Check if the file exists in the local storage
	_, err = os.Stat(repoPath)
	if os.IsNotExist(err) {
		// The file does not exist, nothing to remove
		s.logger.Debug("File does not exist locally, nothing to remove", map[string]interface{}{"localPath": repoPath, "reference": reference.String()})
	} else if err != nil {
		// There was an error checking the file
		s.logger.Error("Error checking file", map[string]interface{}{"error": err, "reference": reference.String()})
		return err
	}

	// Remove the directory
	s.logger.Debug("Removing directory", map[string]interface{}{"localPath": repoPath, "reference": reference.String()})
	err = os.RemoveAll(repoPath)
	if err != nil {
		s.logger.Error("Error removing directory", map[string]interface{}{"error": err, "reference": reference.String()})
		return err
	}

	return nil
}

func (s *service) cloneRepo(ctx context.Context, repoPath string, repoURL string) error {
	_, err := git.PlainCloneContext(ctx, repoPath, false, &git.CloneOptions{
		URL: repoURL,
		// TODO: Fix this
		// Progress: LoggerWriter{s.logger},
	})
	return err
}

func (s *service) pullChanges(ctx context.Context, repoPath string) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	err = w.PullContext(ctx, &git.PullOptions{
		Force: true,
		// Progress: LoggerWriter{s.logger},
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	return nil
}
