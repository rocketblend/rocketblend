package rocketpack

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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

	return &service{
		logger:      options.Logger,
		storagePath: options.StoragePath,
	}, nil
}

func (s *service) GetPackages(ctx context.Context, references ...reference.Reference) (map[reference.Reference]*RocketPack, error) {
	s.logger.Info("Getting packages")
	packages := make(map[reference.Reference]*RocketPack)
	for _, ref := range references {
		s.logger.Info("Processing reference", map[string]interface{}{"reference": ref.String()})

		repo, err := ref.Repo()
		if err != nil {
			s.logger.Error("Error getting repository", map[string]interface{}{"error": err, "reference": ref.String()})
			return nil, err
		}

		repoURL, err := ref.RepoURL()
		if err != nil {
			s.logger.Error("Error getting repository URL", map[string]interface{}{"error": err, "reference": ref.String()})
			return nil, err
		}

		repoPath := filepath.Join(s.storagePath, repo)
		packagePath := filepath.Join(s.storagePath, ref.String(), FileName)

		// Check if the file exists in the local storage
		_, err = os.Stat(packagePath)
		if os.IsNotExist(err) {
			// The file does not exist, clone the repository
			s.logger.Info("File does not exist locally, cloning repository", map[string]interface{}{"repoURL": repoURL, "path": repoPath, "reference": ref.String()})
			_, err = git.PlainCloneContext(ctx, repoPath, false, &git.CloneOptions{
				URL:      repoURL,
				Progress: LoggerWriter{s.logger},
			})
			if err != nil {
				s.logger.Error("Error cloning repository", map[string]interface{}{"error": err, "reference": ref.String()})
				return nil, err
			}
		} else if err != nil {
			// There was an error checking the file
			s.logger.Error("Error checking file", map[string]interface{}{"error": err, "reference": ref.String()})
			return nil, err
		} else {
			// Open the existing repository
			r, err := git.PlainOpen(repoPath)
			if err != nil {
				s.logger.Error("Error opening repository", map[string]interface{}{"error": err, "reference": ref.String()})
				return nil, err
			}

			// Get the working tree
			w, err := r.Worktree()
			if err != nil {
				s.logger.Error("Error getting worktree", map[string]interface{}{"error": err, "reference": ref.String()})
				return nil, err
			}

			// Pull the latest changes from the origin remote and merge into the current branch
			s.logger.Info("Pulling latest changes", map[string]interface{}{"reference": ref.String()})
			err = w.PullContext(ctx, &git.PullOptions{
				Force:    true,
				Progress: LoggerWriter{s.logger},
			})
			if err != nil && err != git.NoErrAlreadyUpToDate {
				s.logger.Error("Error pulling latest changes", map[string]interface{}{"error": err, "reference": ref.String()})
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
			s.logger.Info("Package has dependencies", map[string]interface{}{"reference": ref.String()})

			// Get the dependencies
			depPackages, err := s.GetPackages(ctx, deps...)
			if err != nil {
				s.logger.Error("Error getting dependency packages", map[string]interface{}{"error": err, "reference": ref.String()})
				return nil, err
			}

			// Add the dependencies to the packages map
			for _, dep := range deps {
				packages[dep] = depPackages[dep]
			}

			s.logger.Info("Dependency packages successfully loaded", map[string]interface{}{"reference": ref.String()})
		}

		packages[ref] = pack
	}

	s.logger.Info("Packages successfully loaded")
	return packages, nil
}

func (s *service) RemovePackages(ctx context.Context, references ...reference.Reference) error {
	s.logger.Info("Removing packages")
	for _, ref := range references {
		s.logger.Info("Processing reference", map[string]interface{}{"reference": ref.String()})

		repoPath, err := ref.RepoPath()
		if err != nil {
			s.logger.Error("Error getting repository path", map[string]interface{}{"error": err, "reference": ref.String()})
			return err
		}

		localPath := filepath.Join(s.storagePath, repoPath)

		// Check if the file exists in the local storage
		_, err = os.Stat(localPath)
		if os.IsNotExist(err) {
			// The file does not exist, nothing to remove
			s.logger.Debug("File does not exist locally, nothing to remove", map[string]interface{}{"localPath": localPath, "reference": ref.String()})
			continue
		} else if err != nil {
			// There was an error checking the file
			s.logger.Error("Error checking file", map[string]interface{}{"error": err, "reference": ref.String()})
			return err
		}

		// Remove the directory
		s.logger.Debug("Removing directory", map[string]interface{}{"localPath": localPath, "reference": ref.String()})
		err = os.RemoveAll(localPath)
		if err != nil {
			s.logger.Error("Error removing directory", map[string]interface{}{"error": err, "reference": ref.String()})
			return err
		}
	}

	s.logger.Info("Packages successfully removed")
	return nil
}
