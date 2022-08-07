package local

import (
	"fmt"

	"github.com/spf13/cobra"
)

// installCmdCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "add",
	Short: "Installs a new verison of blender into the local repository",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add called")
	},
}

func init() {
	LocalCmd.AddCommand(installCmd)
}
