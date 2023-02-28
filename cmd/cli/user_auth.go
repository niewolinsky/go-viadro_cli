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

var UserAuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "",
	Long:  ``,
	Run:   userAuth,
	Args:  cobra.ExactArgs(2),
}

func userAuth(cli *cobra.Command, args []string) {
	URL := viper.GetString("endpoint") + "/users/authenticate"

	client := &http.Client{Timeout: 10 * time.Second}

	input := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    args[0],
		Password: args[1],
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

	if res.StatusCode == http.StatusCreated {
		respStruct := struct {
			AuthenticationToken struct {
				Token  string    `json:"token"`
				Expiry time.Time `json:"expiry"`
			} `json:"authentication_token"`
		}{}

		err = json.NewDecoder(res.Body).Decode(&respStruct)
		if err != nil {
			fmt.Println(err)
		}

		viper.Set("tkn", respStruct.AuthenticationToken.Token)
		err = viper.WriteConfig()
		if err != nil {
			panic("Could not write config: " + err.Error())
		}

		fmt.Println("Sucessfully authenticated with token: ", respStruct.AuthenticationToken.Token)
	} else {
		fmt.Println(res.StatusCode)
	}
}

func init() {
	UserCmd.AddCommand(UserAuthCmd)
}
