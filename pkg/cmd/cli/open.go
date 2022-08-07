package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	filePath    string
	blenderArgs string
)

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Opens a rocketfile",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("open called")
	},
}

func init() {
	rootCmd.AddCommand(openCmd)

	openCmd.Flags().StringVarP(&filePath, "file", "f", ".rocket", "The path to the rocketfile to open")
	openCmd.Flags().StringVarP(&blenderArgs, "args", "a", "", "Arguments to pass to blender")

	// if err := openCmd.MarkFlagRequired("file"); err != nil {
	// 	fmt.Println(err)
	// }

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// openCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// openCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
