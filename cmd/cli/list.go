package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Long:  ``,
	Run:   subAlert,
}

func subAlert(cli *cobra.Command, args []string) {
	fmt.Println("test")
}
