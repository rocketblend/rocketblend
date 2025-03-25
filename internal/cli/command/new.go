package command

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rocketblend/rocketblend/internal/cli/ui"
	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/spf13/cobra"
)

type createProjectOpts struct {
	commandOpts
	Name         string
	Overwrite    bool
	ProgressChan chan<- ui.ProgressEvent
}

// newNewCommand creates a new cobra.Command for creating a new project.
// It expects an optional argument which is the name of the project.
func newNewCommand(opts commandOpts) *cobra.Command {
	var overwrite bool
	var name string

	cc := &cobra.Command{
		Use:   "new [name]",
		Short: "Create a new project",
		Long:  `Creates a new project with a specified name.`,
		Args:  cobra.MaximumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// Default project name is derived from the current working directory.
			name = generateProjectName(filepath.Base(opts.Global.WorkingDirectory))
			if len(args) > 0 {
				name = args[0]
			}
			return validateProjectName(name)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			projOpts := createProjectOpts{
				commandOpts: opts,
				Name:        name,
				Overwrite:   overwrite,
			}
			// Use verbose mode directly without UI.
			if opts.Global.Verbose {
				return createProject(cmd.Context(), projOpts)
			}
			return createProjectWithUI(cmd.Context(), projOpts)
		},
	}

	cc.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "overwrite the project if one already exists")
	return cc
}

// createProjectWithUI runs createProject asynchronously and displays progress via the generic progress UI.
func createProjectWithUI(ctx context.Context, opts createProjectOpts) error {
	eventChan := make(chan ui.ProgressEvent, 10)
	opts.ProgressChan = eventChan
	projCtx, cancelProj := context.WithCancel(ctx)
	defer cancelProj()

	go func() {
		defer close(eventChan)
		if err := createProject(projCtx, opts); err != nil {
			eventChan <- ui.ErrorEvent{Message: err.Error()}
		}
	}()

	m := ui.NewProgressModel(eventChan, cancelProj)
	program := tea.NewProgram(&m, tea.WithContext(ctx))
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("failed to run UI: %w", err)
	}

	return nil
}

// createProject creates a new project with the specified name and emits progress events.
func createProject(ctx context.Context, opts createProjectOpts) error {
	emit := func(ev ui.ProgressEvent) {
		if opts.ProgressChan != nil {
			opts.ProgressChan <- ev
		}
	}

	// Check if a project already exists.
	if existingProject(opts.Global.WorkingDirectory) && !opts.Overwrite {
		return errors.New("project already exists in directory")
	}

	emit(ui.StepEvent{Message: "Initializing..."})
	container, err := getContainer(containerOpts{
		AppName:     opts.AppName,
		Development: opts.Development,
		Level:       opts.Global.Level,
		Verbose:     opts.Global.Verbose,
	})
	if err != nil {
		return err
	}

	configurator, err := container.GetConfigurator()
	if err != nil {
		return err
	}

	config, err := configurator.Get()
	if err != nil {
		return err
	}

	driver, err := container.GetDriver()
	if err != nil {
		return err
	}

	emit(ui.StepEvent{Message: "Creating profiles..."})
	profiles, err := driver.LoadProfiles(ctx, &types.LoadProfilesOpts{
		Paths: []string{opts.Global.WorkingDirectory},
		Default: &types.Profile{
			Dependencies: []*types.Dependency{
				{
					Reference: config.DefaultBuild,
					Type:      types.PackageBuild,
				},
			},
		},
	})
	if err != nil {
		return err
	}

	if err := driver.TidyProfiles(ctx, &types.TidyProfilesOpts{
		Profiles: profiles.Profiles,
	}); err != nil {
		return err
	}

	emit(ui.StepEvent{Message: "Installing dependencies..."})
	if err := driver.InstallProfiles(ctx, &types.InstallProfilesOpts{
		Profiles: profiles.Profiles,
	}); err != nil {
		return err
	}

	emit(ui.StepEvent{Message: "Resolving dependencies..."})
	resolveResults, err := driver.ResolveProfiles(ctx, &types.ResolveProfilesOpts{
		Profiles: profiles.Profiles,
	})
	if err != nil {
		return err
	}

	emit(ui.StepEvent{Message: "Creating project..."})
	blendFilePath := filepath.Join(opts.Global.WorkingDirectory, ensureBlendExtension(opts.Name))
	blender, err := container.GetBlender()
	if err != nil {
		return err
	}

	if err := blender.Create(ctx, &types.CreateOpts{
		BlenderOpts: types.BlenderOpts{
			BlendFile: &types.BlendFile{
				Path:         blendFilePath,
				Dependencies: resolveResults.Installations[0],
				Strict:       profiles.Profiles[0].Strict,
			},
			Background: true,
		},
		Overwrite: opts.Overwrite,
	}); err != nil {
		if !errors.Is(err, types.ErrFileExists) {
			return err
		}
	}

	emit(ui.StepEvent{Message: "Saving profiles..."})
	if err := driver.SaveProfiles(ctx, &types.SaveProfilesOpts{
		Profiles: map[string]*types.Profile{
			filepath.Dir(blendFilePath): profiles.Profiles[0],
		},
		EnsurePaths: true,
		Overwrite:   true,
	}); err != nil {
		return err
	}

	emit(ui.CompletionEvent{Message: "Project created!"})
	return nil
}

// existingProject checks if a project already exists at the specified path.
func existingProject(path string) bool {
	return existingBlendFile(path) && existingProfileDir(path)
}

// existingBlendFile checks if a blend file already exists at the specified path.
func existingBlendFile(path string) bool {
	_, err := findFilePathForExt(path, types.BlendFileExtension)
	return err == nil
}

// existingProfileDir checks if a profile directory exists at the specified path.
func existingProfileDir(path string) bool {
	profilePath := filepath.Join(path, types.ProfileDirName)
	info, err := os.Stat(profilePath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// validateProjectName checks if the project name is valid.
func validateProjectName(projectName string) error {
	if filepath.IsAbs(projectName) || strings.Contains(projectName, string(filepath.Separator)) {
		return fmt.Errorf("%q is not a valid project name, it should not contain any path separators", projectName)
	}
	if ext := filepath.Ext(projectName); ext != "" {
		return fmt.Errorf("%q is not a valid project name, it should not contain any file extension", projectName)
	}
	return nil
}

// ensureBlendExtension adds ".blend" extension to filename if it does not already have it.
func ensureBlendExtension(filename string) string {
	if !strings.HasSuffix(filename, ".blend") {
		filename += ".blend"
	}
	return filename
}

// generateProjectName creates a project name by lowercasing and replacing spaces with hyphens.
func generateProjectName(folderName string) string {
	projectName := strings.ToLower(folderName)
	projectName = strings.ReplaceAll(projectName, " ", "-")
	return projectName
}
