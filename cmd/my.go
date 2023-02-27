/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// myCmd represents the my command
var myCmd = &cobra.Command{
	Use:   "my",
	Short: "",
	Long:  ``,
	Run:   listMy,
}

func listMy(cmd *cobra.Command, args []string) {
	URL := "http://localhost:4000/v1/documents/my"

	fmt.Println("Trying to get your documents...")

	// bearer := "Bearer " + os.Getenv("viadro_auth_token")
	bearer := "Bearer " + "TPKQVHZIBQ3YPSWGRAORW6NUFY"
	fmt.Println(bearer)

	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
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
	listCmd.AddCommand(myCmd)

	myCmd.PersistentFlags().StringP("visibility", "v", "", "Possible values: public, hidden")
	myCmd.MarkPersistentFlagRequired("visibility")
}
