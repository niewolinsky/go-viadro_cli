package main

import (
	"fmt"
	"os"
	"viadro_cli/cmd/cli"
	"viadro_cli/config"

	"github.com/spf13/cobra"
)

func main() {
	config.Init()
	var rootCmd = &cobra.Command{Use: "viadro", Short: "A CLI for Viadro API."}
	rootCmd.AddCommand(cli.ListCmd)
	rootCmd.AddCommand(cli.UserCmd)
	rootCmd.AddCommand(cli.UploadCmd)
	rootCmd.AddCommand(cli.MergeCmd)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not execute command: %s\n", err.Error())
		os.Exit(1)
	}
}
