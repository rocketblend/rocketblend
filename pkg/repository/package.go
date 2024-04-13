package repository

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
	"github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/taskrunner"
	"github.com/rocketblend/rocketblend/pkg/types"
)

type (
	getPackageResult struct {
		packs map[reference.Reference]*types.Package
		error error
	}
)

func (r *repository) GetPackages(ctx context.Context, opts *types.GetPackagesOpts) (*types.GetPackagesResult, error) {
	if err := r.validator.Validate(opts); err != nil {
		return nil, err
	}

	packs, err := r.getPackages(ctx, opts.References, opts.Update)
	if err != nil {
		return nil, err
	}

	return &types.GetPackagesResult{
		Packs: packs,
	}, nil
}

func (r *repository) RemovePackages(ctx context.Context, opts *types.RemovePackagesOpts) error {
	if err := r.validator.Validate(opts); err != nil {
		return err
	}

	if err := r.removePackages(ctx, opts.References); err != nil {
		return err
	}

	return nil
}

func (r *repository) InsertPackages(ctx context.Context, opts *types.InsertPackagesOpts) error {
	if err := r.validator.Validate(opts); err != nil {
		return err
	}

	if err := r.insertPackages(ctx, opts.Packs); err != nil {
		return err
	}

	return nil
}

func (r *repository) getPackages(ctx context.Context, references []reference.Reference, update bool) (map[reference.Reference]*types.Package, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make(chan getPackageResult)
	var wg sync.WaitGroup
	wg.Add(len(references))

	for _, ref := range references {
		go func(ref reference.Reference) {
			defer wg.Done()

			packs, err := r.getPackage(ctx, ref, update)
			if err != nil {
				cancel()
				results <- getPackageResult{packs: nil, error: err}
				return
			}

			results <- getPackageResult{packs: packs, error: nil}
		}(ref)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	packages := make(map[reference.Reference]*types.Package)
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

func (r *repository) removePackages(ctx context.Context, references []reference.Reference) error {
	tasks := make([]taskrunner.Task, len(references))
	for _, ref := range references {
		tasks = append(tasks, func(ctx context.Context) error {
			return r.removePackage(ctx, ref)
		})
	}

	return taskrunner.Run(ctx, &taskrunner.RunOpts{
		Tasks: tasks,
		Mode:  taskrunner.Concurrent,
	})
}

func (r *repository) insertPackages(ctx context.Context, packs map[reference.Reference]*types.Package) error {
	tasks := make([]taskrunner.Task, len(packs))
	for ref, pack := range packs {
		tasks = append(tasks, func(ctx context.Context) error {
			return r.insertPackage(ctx, ref, pack)
		})
	}

	return taskrunner.Run(ctx, &taskrunner.RunOpts{
		Tasks: tasks,
		Mode:  taskrunner.Concurrent,
	})
}

func (r *repository) insertPackage(ctx context.Context, ref reference.Reference, pack *types.Package) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	packagePath := filepath.Join(r.storagePath, ref.String(), types.PackageFileName)

	err := os.MkdirAll(filepath.Dir(packagePath), 0755)
	if err != nil {
		r.logger.Error("error creating directory", map[string]interface{}{"error": err, "reference": ref.String(), "path": filepath.Dir(packagePath)})
		return err
	}

	if err := helpers.Save(r.validator, packagePath, pack); err != nil {
		r.logger.Error("error saving package", map[string]interface{}{
			"error":     err,
			"reference": ref.String(),
			"path":      packagePath,
		})

		return err
	}

	return nil
}

func (s *repository) getPackage(ctx context.Context, ref reference.Reference, update bool) (map[reference.Reference]*types.Package, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.logger.Info("processing reference", map[string]interface{}{"reference": ref.String()})

	packages := make(map[reference.Reference]*types.Package)
	repo, err := ref.GetRepo()
	if err != nil {
		s.logger.Error("error getting repository", map[string]interface{}{"error": err, "reference": ref.String()})
		return nil, err
	}

	repoPath := filepath.Join(s.storagePath, repo)
	packagePath := filepath.Join(s.storagePath, ref.String(), types.PackageFileName)

	// The repository does not exist locally, clone it
	if _, err := os.Stat(repoPath); os.IsNotExist(err) || update && !ref.IsLocalOnly() {
		repoURL, err := ref.GetRepoURL()
		if err != nil {
			s.logger.Error("error getting repository URL", map[string]interface{}{"error": err, "reference": ref.String()})
			return nil, err
		}

		s.logger.Info("cloning repository", map[string]interface{}{"repoURL": repoURL, "path": repoPath, "reference": ref.String()})
		if err := s.cloneRepo(ctx, repoPath, repoURL); err != nil {
			return nil, err
		}
	}

	// Check if the file exists in the repository
	if _, err := os.Stat(packagePath); os.IsNotExist(err) || update && !ref.IsLocalOnly() {
		// The file does not exist or forced update, pull the latest changes
		s.logger.Info("pulling latest changes for repository", map[string]interface{}{"path": packagePath, "reference": ref.String()})
		if err := s.pullChanges(ctx, repoPath); err != nil {
			return nil, err
		}
	}

	pack, err := helpers.Load[types.Package](s.validator, packagePath)
	if err != nil {
		s.logger.Error("error loading package", map[string]interface{}{
			"error":     err,
			"reference": ref.String(),
			"path":      packagePath,
		})

		return nil, err
	}

	packages[ref] = pack

	return packages, nil
}

func (s *repository) removePackage(ctx context.Context, reference reference.Reference) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.logger.Info("processing reference", map[string]interface{}{"reference": reference.String()})

	repo, err := reference.GetRepo()
	if err != nil {
		s.logger.Error("error getting repository path", map[string]interface{}{"error": err, "reference": reference.String()})
		return err
	}

	repoPath := filepath.Join(s.storagePath, repo)

	// Check if the file exists in the local storage
	_, err = os.Stat(repoPath)
	if os.IsNotExist(err) {
		// The file does not exist, nothing to remove
		s.logger.Debug("file does not exist locally, nothing to remove", map[string]interface{}{"localPath": repoPath, "reference": reference.String()})
	} else if err != nil {
		// There was an error checking the file
		s.logger.Error("error checking file", map[string]interface{}{"error": err, "reference": reference.String()})
		return err
	}

	// Remove the directory
	s.logger.Debug("removing directory", map[string]interface{}{"localPath": repoPath, "reference": reference.String()})
	err = os.RemoveAll(repoPath)
	if err != nil {
		s.logger.Error("error removing directory", map[string]interface{}{"error": err, "reference": reference.String()})
		return err
	}

	return nil
}

func (s *repository) cloneRepo(ctx context.Context, repoPath string, repoURL string) error {
	_, err := git.PlainCloneContext(ctx, repoPath, false, &git.CloneOptions{
		URL: repoURL,
		// TODO: Fix this
		// Progress: LoggerWriter{s.logger},
	})
	return err
}

func (s *repository) pullChanges(ctx context.Context, repoPath string) error {
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
