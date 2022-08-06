/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// verisonCmd represents the verison command
var verisonCmd = &cobra.Command{
	Use:   "verison",
	Short: "Prints RockerBlends version number",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("verison called")
	},
}

func init() {
	rootCmd.AddCommand(verisonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verisonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// verisonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
