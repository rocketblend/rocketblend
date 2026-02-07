package supervisor

import (
	"context"
	"sync"

	"github.com/rocketblend/rocketblend/pkg/types"
)

type ProjectStatus string

const (
	StatusRegistered ProjectStatus = "registered"
	StatusActive     ProjectStatus = "active"
	StatusInactive   ProjectStatus = "inactive"
)

type (
	Project struct {
		Path      string
		BlendFile string
		Status    ProjectStatus
		Profile   *types.Profile
	}

	Supervisor struct {
		projects map[string]*Project
		mutex    sync.RWMutex
		logger   types.Logger
		driver   types.Driver
	}
)

func New(logger types.Logger, driver types.Driver) *Supervisor {
	return &Supervisor{
		projects: make(map[string]*Project),
		logger:   logger,
		driver:   driver,
	}
}

func (s *Supervisor) Discover(ctx context.Context, root string) error {
	projectFiles, err := Discover(ctx, root)
	if err != nil {
		return err
	}

	if len(projectFiles) == 0 {
		s.logger.Info("no projects discovered", map[string]any{"root": root})
		return nil
	}

	projectDirs := make([]string, 0, len(projectFiles))
	for dir := range projectFiles {
		projectDirs = append(projectDirs, dir)
	}

	result, err := s.driver.LoadProfiles(ctx, &types.LoadProfilesOpts{
		Paths: projectDirs,
	})
	if err != nil {
		s.logger.Error("failed to load profiles", map[string]any{"err": err})
		return err
	}

	for i, profile := range result.Profiles {
		dir := projectDirs[i]
		blendFile := projectFiles[dir]
		s.register(dir, blendFile, profile)
	}

	s.logger.Info("discovery complete", map[string]any{
		"discovered": len(projectFiles),
		"registered": len(result.Profiles),
	})

	return nil
}

func (s *Supervisor) Registered() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	paths := make([]string, 0, len(s.projects))
	for path := range s.projects {
		paths = append(paths, path)
	}

	return paths
}

func (s *Supervisor) register(path string, blendFile string, profile *types.Profile) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.projects[path]; exists {
		return
	}

	s.projects[path] = &Project{
		Path:      path,
		BlendFile: blendFile,
		Status:    StatusRegistered,
		Profile:   profile,
	}

	s.logger.Info("registered project", map[string]any{"path": path, "blend_file": blendFile})
}
