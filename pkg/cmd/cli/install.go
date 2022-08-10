package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/rocketblend/rocketblend/pkg/client/installer"
	"github.com/spf13/cobra"
)

var (
	buildHash       string
	installationDir string
)

// installCmdCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs a new verison of blender into the local repository",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Installation started")
		build := client.FindAvailableBuildFromHash(buildHash)
		if build == nil {
			fmt.Printf("No build found with hash %s\n", buildHash)
			return
		}

		fileName := filepath.Join(installationDir, build.Name)
		err := installer.Install(string(build.Uri), fileName)
		if err != nil {
			fmt.Printf("Error installing build: %s\n", err)
		}

		// Possibly strip extension from build name on remote.
		name := strings.TrimSuffix(build.Name, filepath.Ext(build.Name))
		install := client.Install{
			Name:    name,
			Hash:    build.Hash,
			Path:    filepath.Join(installationDir, name),
			Version: build.Version.String(),
		}

		client.AddInstall(install)

		fmt.Println("Installation complete!")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVarP(&buildHash, "build", "b", "", "Build hash of the version to install")
	installCmd.Flags().StringVarP(&installationDir, "path", "p", client.GetInstallationDir(), "Path to the installation directory")
	installCmd.MarkFlagRequired("build")
}
