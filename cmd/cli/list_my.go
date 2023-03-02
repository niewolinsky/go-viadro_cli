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

var myCmd = &cobra.Command{
	Use:   "my",
	Short: "List all of user's documents",
	Long:  ``,
	Run:   listMy,
}

func listMy(cli *cobra.Command, args []string) {
	URL := "http://localhost:4000/v1/documents/my"
	bearer := "Bearer " + viper.GetString("tkn")

	visibility, err := cli.Flags().GetString("visibility")
	if err != nil {
		log.Fatal(err)
	}
	//TODO: FLAG DIFFERENT RESULTS
	fmt.Println(visibility)

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal("Can't form request")
	}
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
				Is_hidden  bool      `json:"is_hidden"`
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
			if visibility == "hidden" && !document.Is_hidden {
				continue
			}

			if visibility == "public" && document.Is_hidden {
				continue
			}

			fmt.Println(document)
		}
	} else if res.StatusCode == http.StatusUnauthorized {
		fmt.Println("Invalid or expired token, use auth command to grab a new token.")
	} else {
		fmt.Println("internal server error, try again later", res.StatusCode)
	}
}

func init() {
	ListCmd.AddCommand(myCmd)

	myCmd.PersistentFlags().StringP("visibility", "v", "all", "Possible values: all, public, hidden")
}
