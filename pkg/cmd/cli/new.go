package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	projectName string
	path        string
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new rocketfile",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("new called")
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&projectName, "name", "n", "", "Name of the project")
	if err := newCmd.MarkFlagRequired("name"); err != nil {
		fmt.Println(err)
	}

	newCmd.Flags().StringVarP(&path, "path", "p", "", "Path to create the project at")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
