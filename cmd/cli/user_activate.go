package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var UserActivateCmd = &cobra.Command{
	Use:   "activate",
	Short: "",
	Long:  ``,
	Run:   userActivate,
	Args:  cobra.ExactArgs(1),
}

func userActivate(cli *cobra.Command, args []string) {
	URL := viper.GetString("endpoint") + "/users/activate"

	client := &http.Client{Timeout: 10 * time.Second}

	input := struct {
		TokenPlaintext string `json:"token"`
	}{
		TokenPlaintext: args[0],
	}

	jsonified, _ := json.MarshalIndent(input, "", "\t")
	reader := bytes.NewReader(jsonified)

	req, err := http.NewRequest(http.MethodPut, URL, reader)
	if err != nil {
		fmt.Println("error tutaj")
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(res.StatusCode)
	}
}

func init() {
	UserCmd.AddCommand(UserActivateCmd)
}
