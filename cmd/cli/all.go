package cli

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "",
	Long:  ``,
	Run:   listAll,
}

func listAll(cli *cobra.Command, args []string) {
	URL := "http://localhost:4000/v1/documents"

	fmt.Println("Trying to get all documents...")

	response, err := http.Get(URL)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	} else {
		fmt.Println("Error getting all documents")
	}
}

func init() {
	ListCmd.AddCommand(allCmd)
}
