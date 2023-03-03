package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var DocumentListCmd = &cobra.Command{
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
		respStruct := struct {
			Documents []struct {
				DocumentID int       `json:"document_id"`
				Title      string    `json:"title"`
				Link       string    `json:"link"`
				Tags       []string  `json:"tags"`
				CreatedAt  time.Time `json:"created_at"`
			} `json:"documents"`
			Metadata struct {
				CurrentPage  int `json:"current_page"`
				PageSize     int `json:"page_size"`
				FirstPage    int `json:"first_page"`
				LastPage     int `json:"last_page"`
				TotalRecords int `json:"total_records"`
			} `json:"metadata"`
		}{}

		err = json.NewDecoder(res.Body).Decode(&respStruct)
		if err != nil {
			fmt.Println(err)
		}

		for _, document := range respStruct.Documents {
			fmt.Println(document)
		}

	} else {
		fmt.Println("internal server error, try again later", res.StatusCode)
	}
}

func init() {
	DocumentCmd.AddCommand(DocumentListCmd)
}
