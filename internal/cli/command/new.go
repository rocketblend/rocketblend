package command

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/types"
	"github.com/spf13/cobra"
)

type (
	createProjectOpts struct {
		commandOpts
		Name      string
		Overwrite bool
	}
)

// newNewCommand creates a new cobra.Command object initialized for creating a new project.
// It expects a single argument which is the name of the project.
func newNewCommand(opts commandOpts) *cobra.Command {
	var overwrite bool

	cc := &cobra.Command{
		Use:   "new [name]",
		Short: "Create a new project",
		Long:  `Creates a new project with a specified name.`,
		Args:  cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateProjectName(args[0])
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWithSpinner(cmd.Context(), func(ctx context.Context) error {
				if err := createProject(ctx, createProjectOpts{
					commandOpts: opts,
					Name:        args[0],
					Overwrite:   overwrite,
				}); err != nil {
					return fmt.Errorf("failed to create project: %w", err)
				}

				return nil
			}, &spinnerOptions{
				Suffix:  "Creating project...",
				Verbose: opts.Global.Verbose,
			})
		},
	}

	cc.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "overwrite the project if one already exists")

	return cc
}

// createProject creates a new project with the specified name.
func createProject(ctx context.Context, opts createProjectOpts) error {
	if !opts.Overwrite && existingProject(opts.Global.WorkingDirectory) {
		return errors.New("project already exists in directory")
	}

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

	profiles := []*types.Profile{
		{
			Dependencies: []*types.Dependency{
				{
					Reference: config.DefaultBuild,
					Type:      types.PackageBuild,
				},
			},
			InjectionMode: types.DefaultInjectionMode,
		},
	}

	if err := driver.TidyProfiles(ctx, &types.TidyProfilesOpts{
		Profiles: profiles,
	}); err != nil {
		return err
	}

	if err := driver.InstallProfiles(ctx, &types.InstallProfilesOpts{
		Profiles: profiles,
	}); err != nil {
		return err
	}

	resolveResults, err := driver.ResolveProfiles(ctx, &types.ResolveProfilesOpts{
		Profiles: profiles,
	})
	if err != nil {
		return err
	}

	blender, err := container.GetBlender()
	if err != nil {
		return err
	}

	blendFilePath := filepath.Join(opts.Global.WorkingDirectory, ensureBlendExtension(opts.Name))

	if err := blender.Create(ctx, &types.CreateOpts{
		BlenderOpts: types.BlenderOpts{
			BlendFile: &types.BlendFile{
				Path:          blendFilePath,
				Dependencies:  resolveResults.Installations[0],
				InjectionMode: profiles[0].InjectionMode,
			},
			Background: true,
		},
	}); err != nil {
		return err
	}

	if err := driver.SaveProfiles(ctx, &types.SaveProfilesOpts{
		Profiles: map[string]*types.Profile{
			filepath.Dir(blendFilePath): profiles[0],
		},
	}); err != nil {
		return err
	}

	return nil
}

// existingProject checks if a project already exists at the specified path.
func existingProject(path string) bool {
	_, err := findFilePathForExt(path, types.BlendFileExtension)
	return err == nil
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
