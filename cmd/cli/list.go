package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Long:  ``,
	Run: func(cli *cobra.Command, args []string) {
		fmt.Println("Available subcommands: all, my")
	},
}
