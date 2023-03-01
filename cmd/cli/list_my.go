package cli

import (
	"fmt"
	"io"
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
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
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
