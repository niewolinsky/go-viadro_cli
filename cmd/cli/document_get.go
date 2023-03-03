package cli

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var DocumentGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details about a document",
	Long:  ``,
	Run:   documentGet,
	Args:  cobra.ExactArgs(1),
}

func documentGet(cli *cobra.Command, args []string) {
	document_id, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf(`http://localhost:4000/v1/document/%d`, document_id)
	bearer := "Bearer " + viper.GetString("tkn")

	_, err = cli.Flags().GetBool("details")
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal("Service unavailable, try again later.")
	}

	req.Header.Add("Authorization", bearer)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Service unavailable, try again later.")
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		//? destruct it and grab each fragment, then use flag to extract only url or all details
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	} else if res.StatusCode == http.StatusUnauthorized {
		fmt.Println("You do not have permissions to view the document.")
	} else if res.StatusCode == http.StatusNotFound {
		fmt.Println("Document with given ID does not exist.")
	} else {
		fmt.Println("Internal server error, try again later.", res.StatusCode)
	}
}

func init() {
	DocumentCmd.AddCommand(DocumentGetCmd)
	DocumentGetCmd.PersistentFlags().Bool("details", false, "See file details? Default: hidden")
}
