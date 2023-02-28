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

var UserRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "",
	Long:  ``,
	Run:   userRegister,
	Args:  cobra.ExactArgs(3),
}

func userRegister(cli *cobra.Command, args []string) {
	URL := viper.GetString("endpoint") + "/users/register"

	client := &http.Client{Timeout: 10 * time.Second}

	input := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Username: args[0],
		Email:    args[1],
		Password: args[2],
	}

	jsonified, _ := json.MarshalIndent(input, "", "\t")
	reader := bytes.NewReader(jsonified)

	req, err := http.NewRequest(http.MethodPost, URL, reader)
	if err != nil {
		fmt.Println("error tutaj")
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusAccepted {
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
	UserCmd.AddCommand(UserRegisterCmd)
}
