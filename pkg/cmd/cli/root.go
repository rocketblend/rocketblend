package cli

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/rocketblend/rocketblend/pkg/cmd/cli/remote"
)

const verison = "0.0.1"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "rocketblend-cli",
	Version: verison,
	Short:   "Verison manager for blender.",
	Long: `RocketBlend-Cli is a CLI tool for Blender manages verisons of blender,
allowing users to quickly switch between verisons.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addSubcommandPalettes() {
	rootCmd.AddCommand(remote.RemoteCmd)
}

func init() {
	rootCmd.SetVersionTemplate("{{.Version}}\n")

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rocketblend.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	addSubcommandPalettes()
}
