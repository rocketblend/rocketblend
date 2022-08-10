package cli

import (
	"fmt"
	"log"
	"sort"

	"github.com/blang/semver/v4"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/spf13/cobra"
)

// var (
// 	includeAll   bool
// 	notInstalled bool
// 	filterTag    string
// )

type Available struct {
	Hash    string
	Name    string
	Path    string
	Version semver.Version
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the versions of blender that are available",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		builds, err := client.FetchAvailableBuilds()
		if err != nil {
			log.Fatal(err)
		}

		available := []Available{}
		installed := client.GetInstalledBuilds()
		for _, installed := range installed {
			temp := Available{}
			temp.Hash = installed.Hash
			temp.Name = installed.Name
			temp.Version, _ = semver.Parse(installed.Version)
			temp.Path = installed.Path
			available = append(available, temp)
		}

		for _, build := range builds {
			isExisting := false
			for _, existing := range available {
				if build.Hash == existing.Hash {
					isExisting = true
					break
				}
			}

			if !isExisting {
				temp := Available{}
				temp.Hash = build.Hash
				temp.Name = build.Name
				temp.Version, _ = semver.Parse(build.Version)
				temp.Path = ""
				available = append(available, temp)
			}
		}

		sort.SliceStable(available, func(i, j int) bool {
			return available[i].Version.GT(available[j].Version)
		})

		outputBuilds(available)
	},
}

func outputBuilds(available []Available) {
	for _, build := range available {
		fmt.Printf("(%s)\t%s\t\t%s\n", build.Hash, build.Name, build.Path)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)

	// listCmd.Flags().StringVarP(&filterTag, "tag", "t", "", "Filter by tag")

	// listCmd.Flags().BoolVarP(&includeAll, "all", "a", false, "All known available blender versions")
	// listCmd.Flags().BoolVarP(&notInstalled, "not-installed", "n", false, "Versions not locally installed")
	// listCmd.MarkFlagsMutuallyExclusive("all", "not-installed")
}
