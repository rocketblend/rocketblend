package repository

import (
	"context"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/reference"
	"github.com/rocketblend/rocketblend/pkg/taskrunner"
	"github.com/rocketblend/rocketblend/pkg/types"
)

type (
	getPackageResult struct {
		Reference reference.Reference
		Package   *types.Package
	}
)

func (r *Repository) GetPackages(ctx context.Context, opts *types.GetPackagesOpts) (*types.GetPackagesResult, error) {
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

func (r *Repository) RemovePackages(ctx context.Context, opts *types.RemovePackagesOpts) error {
	if err := r.validator.Validate(opts); err != nil {
		return err
	}

	if err := r.removePackages(ctx, opts.References); err != nil {
		return err
	}

	return nil
}

func (r *Repository) InsertPackages(ctx context.Context, opts *types.InsertPackagesOpts) error {
	if err := r.validator.Validate(opts); err != nil {
		return err
	}

	if err := r.insertPackages(ctx, opts.Packs); err != nil {
		return err
	}

	return nil
}

func (r *Repository) getPackages(ctx context.Context, references []reference.Reference, update bool) (map[reference.Reference]*types.Package, error) {
	tasks := make([]taskrunner.Task[*getPackageResult], len(references))
	for _, ref := range references {
		tasks = append(tasks, func(ctx context.Context) (*getPackageResult, error) {
			pack, err := r.getPackage(ctx, ref, update)
			if err != nil {
				return nil, err
			}

			return &getPackageResult{Reference: ref, Package: pack}, nil
		})
	}

	results, err := taskrunner.Run(ctx, &taskrunner.RunOpts[*getPackageResult]{
		Tasks: tasks,
		Mode:  taskrunner.Concurrent,
	})
	if err != nil {
		return nil, err
	}

	packages := make(map[reference.Reference]*types.Package)
	for _, res := range results {
		packages[res.Reference] = res.Package
	}

	return packages, nil
}

func (r *Repository) removePackages(ctx context.Context, references []reference.Reference) error {
	tasks := make([]taskrunner.Task[struct{}], len(references))
	for _, ref := range references {
		tasks = append(tasks, func(ctx context.Context) (struct{}, error) {
			return struct{}{}, r.removePackage(ctx, ref)
		})
	}

	_, err := taskrunner.Run(ctx, &taskrunner.RunOpts[struct{}]{
		Tasks: tasks,
		Mode:  taskrunner.Concurrent,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) insertPackages(ctx context.Context, packs map[reference.Reference]*types.Package) error {
	tasks := make([]taskrunner.Task[struct{}], len(packs))
	for ref, pack := range packs {
		tasks = append(tasks, func(ctx context.Context) (struct{}, error) {
			return struct{}{}, r.insertPackage(ctx, ref, pack)
		})
	}

	_, err := taskrunner.Run(ctx, &taskrunner.RunOpts[struct{}]{
		Tasks: tasks,
		Mode:  taskrunner.Concurrent,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) insertPackage(ctx context.Context, ref reference.Reference, pack *types.Package) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	packagePath := filepath.Join(r.packagePath, ref.String(), types.PackageFileName)

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

func (s *Repository) getPackage(ctx context.Context, ref reference.Reference, update bool) (*types.Package, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.logger.Info("processing reference", map[string]interface{}{"reference": ref.String()})

	repo, err := ref.GetRepo()
	if err != nil {
		s.logger.Error("error getting repository", map[string]interface{}{"error": err, "reference": ref.String()})
		return nil, err
	}

	repoPath := filepath.Join(s.packagePath, repo)
	packagePath := filepath.Join(s.packagePath, ref.String(), types.PackageFileName)

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

	return pack, nil
}

func (s *Repository) removePackage(ctx context.Context, reference reference.Reference) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.logger.Info("processing reference", map[string]interface{}{"reference": reference.String()})

	repo, err := reference.GetRepo()
	if err != nil {
		s.logger.Error("error getting repository path", map[string]interface{}{"error": err, "reference": reference.String()})
		return err
	}

	repoPath := filepath.Join(s.packagePath, repo)

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

func (s *Repository) cloneRepo(ctx context.Context, repoPath string, repoURL string) error {
	_, err := git.PlainCloneContext(ctx, repoPath, false, &git.CloneOptions{
		URL: repoURL,
		// TODO: Fix this
		// Progress: LoggerWriter{s.logger},
	})
	return err
}

func (s *Repository) pullChanges(ctx context.Context, repoPath string) error {
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
