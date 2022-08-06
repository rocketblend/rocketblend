package local

import (
	"fmt"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets details of a version of blender in the local repository",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("set called")
	},
}

func init() {
	LocalCmd.AddCommand(setCmd)
}
