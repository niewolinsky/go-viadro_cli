package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"viadro_cli/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var UserAuthCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate (login) an existing user",
	Long:  ``,
	Run:   userAuth,
	Args:  cobra.ExactArgs(2),
}

func userAuth(cli *cobra.Command, args []string) {
	input := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    args[0],
		Password: args[1],
	}

	req := utils.PrepareRequest(input, viper.GetString("endpoint")+"/users/authenticate", http.MethodPut)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Service unavailable, try again later.")
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

		fmt.Println("Successfully logged in, your authentication token: ", respStruct.AuthenticationToken.Token)
		fmt.Println("Token has been stored in config, for the next 24 hours requests will automatically use it.")
	} else if res.StatusCode == http.StatusBadRequest {
		fmt.Println("malformed json request", res.StatusCode)
	} else if res.StatusCode == http.StatusUnauthorized {
		fmt.Println("wrong email or password", res.StatusCode)
	} else {
		fmt.Println("internal server error, try again later", res.StatusCode)
	}
}

func init() {
	UserCmd.AddCommand(UserAuthCmd)
}
