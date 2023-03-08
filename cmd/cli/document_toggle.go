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
	Run:   DocumentToggle,
	Args:  cobra.ExactArgs(1),
}

func DocumentToggle(cli *cobra.Command, args []string) {
	msg := Toggle(args)
	fmt.Println(msg)
}

func Toggle(args []string) string {
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
		return fmt.Sprintf("Successfully toggled visibility of document with ID: %d", document_id)
	} else if res.StatusCode == http.StatusUnauthorized {
		return "You do not have permissions to view the document."
	} else if res.StatusCode == http.StatusNotFound {
		return "Document with given ID does not exist."
	} else {
		return "Internal server error, try again later."
	}
}

func init() {
	DocumentCmd.AddCommand(DocumentToggleCmd)
}
