package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// userCmd represents the user command
var DocumentCmd = &cobra.Command{
	Use:   "document",
	Short: "Manage documents",
	Long:  ``,
	Run: func(cli *cobra.Command, args []string) {
		fmt.Println("Available subcommands: list, get, grab, merge, toggle, upload")
	},
}
