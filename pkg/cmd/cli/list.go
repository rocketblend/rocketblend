package cli

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/spf13/cobra"
)

// var (
// 	includeAll   bool
// 	notInstalled bool
// 	filterTag    string
// )

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the versions of blender that are available",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		outputBuilds(client.GetAvilableBuilds())
	},
}

func outputBuilds(available []client.Available) {
	for _, build := range available {
		fmt.Printf("(%s)\t%s\t\t%s\n", build.Hash, build.Name, build.Uri)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)

	// listCmd.Flags().StringVarP(&filterTag, "tag", "t", "", "Filter by tag")

	// listCmd.Flags().BoolVarP(&includeAll, "all", "a", false, "All known available blender versions")
	// listCmd.Flags().BoolVarP(&notInstalled, "not-installed", "n", false, "Versions not locally installed")
	// listCmd.MarkFlagsMutuallyExclusive("all", "not-installed")
}
