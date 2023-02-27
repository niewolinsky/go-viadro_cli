/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Long:  ``,
	Run:   subAlert,
}

func subAlert(cmd *cobra.Command, args []string) {
	fmt.Println("Error: must also specify a resource 'all' or 'my'.")
}

func init() {
	rootCmd.AddCommand(listCmd)
}
