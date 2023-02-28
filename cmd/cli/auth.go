package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var AuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "A brief description of your command",
	Long:  ``,
	Run:   auth,
	Args:  cobra.ExactArgs(2),
}

func auth(cli *cobra.Command, args []string) {
	URL := "http://localhost:4000/v1/users/authenticate"

	input := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    args[0],
		Password: args[1],
	}

	jsonified, _ := json.MarshalIndent(input, "", "\t")
	reader := bytes.NewReader(jsonified)

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest(http.MethodPut, URL, reader)
	if err != nil {
		fmt.Println("error tutaj")
	}

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("asd")
	}
	defer response.Body.Close()

	respStruct := struct {
		AuthenticationToken struct {
			Token  string    `json:"token"`
			Expiry time.Time `json:"expiry"`
		} `json:"authentication_token"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&respStruct)
	if err != nil {
		fmt.Println(err)
	}

	viper.Set("tkn", respStruct.AuthenticationToken.Token)
	err = viper.WriteConfig()
	if err != nil {
		panic("Could not write config: " + err.Error())
	}

	fmt.Println("Sucessfully authenticated with token: ", respStruct.AuthenticationToken.Token)
}
