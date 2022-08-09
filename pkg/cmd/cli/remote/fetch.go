package remote

import (
	"log"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Pulls details about avaiable versions of blender from the remote repositories",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := client.FetchAvailableBuilds()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Fetch complete")
	},
}

func init() {
	RemoteCmd.AddCommand(fetchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
