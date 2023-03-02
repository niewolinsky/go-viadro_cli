package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// userCmd represents the user command
var AdminDocumentListCmd = &cobra.Command{
	Use:   "list",
	Short: "Manage users and documents as administrator",
	Long:  ``,
	Run:   listAdminAll,
}

func listAdminAll(cmd *cobra.Command, args []string) {
	URL := "http://localhost:4000/v1/documents/all"

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal("Can't form request")
	}

	bearer := "Bearer " + viper.GetString("tkn")
	req.Header.Add("Authorization", bearer)

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
	} else if res.StatusCode == http.StatusUnauthorized {
		fmt.Println("Invalid or expired token, use auth command to grab a new token.")
	} else if res.StatusCode == http.StatusForbidden {
		fmt.Println("you dont have required role or privilege to complete this action", res.StatusCode)
	} else {
		fmt.Println("internal server error, try again later", res.StatusCode)
	}
}

func init() {
	AdminDocumentCmd.AddCommand(AdminDocumentListCmd)
}
