package cli

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "List all visible (public) documents",
	Long:  ``,
	Run:   listAll,
}

func listAll(cli *cobra.Command, args []string) {
	URL := "http://localhost:4000/v1/documents"

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal("Can't form request")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Service unavailable, try again later.")
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	} else {
		fmt.Println("internal server error, try again later", res.StatusCode)
	}
}

func init() {
	ListCmd.AddCommand(allCmd)
}
