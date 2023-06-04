package rocketpack

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/go-git/go-git/v5"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/reference"
)

type (
	Service interface {
		GetPackages(ctx context.Context, references ...reference.Reference) (map[reference.Reference]*RocketPack, error)
		RemovePackages(ctx context.Context, references ...reference.Reference) error
	}

	Options struct {
		Logger      logger.Logger
		StoragePath string
	}

	Option func(*Options)

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

func (s *service) GetPackages(ctx context.Context, references ...reference.Reference) (map[reference.Reference]*RocketPack, error) {
	rpm := NewRocketPackMap()
	errs := make(chan error, len(references))
	var wg sync.WaitGroup
	wg.Add(len(references))

	for _, ref := range references {
		go func(ref reference.Reference) {
			defer wg.Done()
			packs, err := s.getPackages(ctx, ref)
			if err != nil {
				errs <- fmt.Errorf("error getting package %s: %w", ref, err)
				return
			}

			for index, pack := range packs {
				rpm.Store(index, pack)
			}
		}(ref)
	}

	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		return nil, fmt.Errorf("errors occurred: %v", <-errs) // return first error for simplicity
	}

	return rpm.ToRegularMap(), nil
}

func (s *service) RemovePackages(ctx context.Context, references ...reference.Reference) error {
	errs := make(chan error, len(references))
	var wg sync.WaitGroup
	wg.Add(len(references))

	for _, ref := range references {
		go func(ref reference.Reference) {
			defer wg.Done()
			err := s.removePackage(ctx, ref)
			if err != nil {
				errs <- fmt.Errorf("error removing package %s: %w", ref, err)
				return
			}
		}(ref)
	}

	wg.Wait()
	close(errs)

	if len(errs) > 0 {
		return fmt.Errorf("errors occurred: %v", <-errs) // return first error for simplicity
	}

	return nil
}

func (s *service) getPackages(ctx context.Context, ref reference.Reference) (map[reference.Reference]*RocketPack, error) {
	s.logger.Info("Processing reference", map[string]interface{}{"reference": ref.String()})

	packages := make(map[reference.Reference]*RocketPack)

	repo, err := ref.GetRepo()
	if err != nil {
		s.logger.Error("Error getting repository", map[string]interface{}{"error": err, "reference": ref.String()})
		return nil, err
	}

	repoURL, err := ref.GetRepoURL()
	if err != nil {
		s.logger.Error("Error getting repository URL", map[string]interface{}{"error": err, "reference": ref.String()})
		return nil, err
	}

	repoPath := filepath.Join(s.storagePath, repo)
	packagePath := filepath.Join(s.storagePath, ref.String(), FileName)

	// The repository does not exist locally, clone it
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		s.logger.Info("Repository does not exist locally, cloning repository", map[string]interface{}{"repoURL": repoURL, "path": repoPath, "reference": ref.String()})
		if err := s.cloneRepo(ctx, repoPath, repoURL); err != nil {
			return nil, fmt.Errorf("error cloning repository: %w", err)
		}
	}

	// Check if the file exists in the repository
	if _, err := os.Stat(packagePath); os.IsNotExist(err) {
		// The file does not exist, pull the latest changes
		s.logger.Info("File does not exist locally, pulling latest changes", map[string]interface{}{"path": packagePath, "reference": ref.String()})
		if err := s.pullChanges(ctx, repoPath); err != nil {
			return nil, fmt.Errorf("error pulling latest changes: %w", err)
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
		depPackages, err := s.GetPackages(ctx, deps...)
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
		URL:      repoURL,
		Progress: LoggerWriter{s.logger},
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
		Force:    true,
		Progress: LoggerWriter{s.logger},
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	return nil
}
