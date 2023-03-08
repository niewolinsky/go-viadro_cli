package cli

import (
	"log"
	"os"

	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/spf13/cobra"
)

var MergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "",
	Run:   documentsMerge,
}

func documentsMerge(cli *cobra.Command, args []string) {
	myFile, err := os.Create("merge-output.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer myFile.Close()

	err = pdfcpu.Merge("sample3.pdf", []string{"sample1.pdf", "sample2.pdf"}, myFile, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	DocumentCmd.AddCommand(MergeCmd)
}
