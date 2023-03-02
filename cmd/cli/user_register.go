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

var UserRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user for the Viadro service",
	Long:  ``,
	Run:   userRegister,
	Args:  cobra.ExactArgs(3),
}

func userRegister(cli *cobra.Command, args []string) {
	input := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Username: args[0],
		Email:    args[1],
		Password: args[2],
	}

	req := utils.PrepareRequest(input, viper.GetString("endpoint")+"/user/register", http.MethodPost)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Service unavailable, try again later.")
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusAccepted {
		respStruct := struct {
			User struct {
				ID        int       `json:"id"`
				CreatedAt time.Time `json:"created_at"`
				Username  string    `json:"username"`
				Email     string    `json:"email"`
				Activated bool      `json:"activated"`
			} `json:"user"`
		}{}

		err = json.NewDecoder(res.Body).Decode(&respStruct)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Created new user with ID: %d, check your email for account activation. \n", respStruct.User.ID)
	} else if res.StatusCode == http.StatusBadRequest {
		fmt.Println("malformed json request", res.StatusCode)
	} else if res.StatusCode == http.StatusUnprocessableEntity {
		fmt.Println("account with this email already exists", res.StatusCode)
	} else {
		fmt.Println("internal server error, try again later", res.StatusCode)
	}
}

func init() {
	UserCmd.AddCommand(UserRegisterCmd)
}
