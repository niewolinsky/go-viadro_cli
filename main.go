package main

import (
	"viadro_cli/cmd/cli"
	"viadro_cli/cmd/tui"
	"viadro_cli/config"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func main() {
	config.Init()

	rootCmd := &cobra.Command{Use: "viadro", Short: "A CLI tool for Viadro API."}
	rootCmd.AddCommand(cli.DocumentCmd)
	rootCmd.AddCommand(cli.UserCmd)
	rootCmd.AddCommand(cli.AdminCmd)
	rootCmd.AddCommand(tui.TuiCmd)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("invalid command")
	}
}
