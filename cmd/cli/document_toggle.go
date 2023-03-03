package cli

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var DocumentToggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggle document visibility",
	Long:  ``,
	Run:   documentToggle,
	Args:  cobra.ExactArgs(1),
}

func documentToggle(cli *cobra.Command, args []string) {
	document_id, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf(`http://localhost:4000/v1/document/%d`, document_id)
	bearer := "Bearer " + viper.GetString("tkn")

	req, err := http.NewRequest(http.MethodPatch, url, nil)
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
		fmt.Println("Successfully toggled visibility of document with ID:", document_id)
	} else if res.StatusCode == http.StatusUnauthorized {
		fmt.Println("You do not have permissions to view the document.")
	} else if res.StatusCode == http.StatusNotFound {
		fmt.Println("Document with given ID does not exist.")
	} else {
		fmt.Println("Internal server error, try again later.", res.StatusCode)
	}

}

func init() {
	DocumentCmd.AddCommand(DocumentToggleCmd)
}
