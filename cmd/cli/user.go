package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"viadro_cli/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// * COMMANDS * //
var UserCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage user's credentials",
	Run: func(cli *cobra.Command, args []string) {
		fmt.Println("Available subcommands: register, activate, auth, delete")
	},
}
var UserRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user for the Viadro service",
	Run:   userRegister,
	Args:  cobra.ExactArgs(3),
}
var UserActivateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Activate user account",
	Run:   userActivate,
	Args:  cobra.ExactArgs(1),
}
var UserAuthCmd = &cobra.Command{
	Use:     "auth",
	Aliases: []string{"login"},
	Short:   "Authenticate (login) an existing user",
	Run:     userAuth,
	Args:    cobra.ExactArgs(2),
}
var UserDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete (deactivate) a user from Viadro service",
	Run:   userDelete,
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

	switch res.StatusCode {
	case http.StatusAccepted:
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

		fmt.Printf("created new user with ID: %d, check your email for account activation. \n", respStruct.User.ID)
	case http.StatusBadRequest:
		fmt.Println("malformed json request")
	case http.StatusUnprocessableEntity:
		fmt.Println("account with this email already exists")
	default:
		fmt.Println("internal server error, try again later")
	}
}

func userActivate(cli *cobra.Command, args []string) {
	input := struct {
		TokenPlaintext string `json:"token"`
	}{
		TokenPlaintext: args[0],
	}

	req := utils.PrepareRequest(input, viper.GetString("endpoint")+"/user/activate", http.MethodPut)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Service unavailable, try again later.")
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		fmt.Println("User successfully activated, use auth command to login.")
	case http.StatusBadRequest:
		fmt.Println("malformed json request")
	case http.StatusUnprocessableEntity:
		fmt.Println("invalid or expired token")
	default:
		fmt.Println("internal server error, try again later", res.StatusCode)
	}
}

func userAuth(cli *cobra.Command, args []string) {
	input := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    args[0],
		Password: args[1],
	}

	req := utils.PrepareRequest(input, viper.GetString("endpoint")+"/user/authenticate", http.MethodPut)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Service unavailable, try again later.")
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusCreated:
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
	case http.StatusBadRequest:
		fmt.Println("malformed json request")
	case http.StatusUnauthorized:
		fmt.Println("wrong email or password")
	default:
		fmt.Println("internal server error, try again later")
	}
}

func userDelete(cli *cobra.Command, args []string) {
	user_id, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf(`http://localhost:4000/v1/user/%d`, user_id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Fatal("Service unavailable, try again later.")
	}

	bearer := "Bearer " + viper.GetString("tkn")
	req.Header.Add("Authorization", bearer)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Service unavailable, try again later.")
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusNoContent:
		fmt.Println("User successfully deleted.")
	case http.StatusUnauthorized:
		fmt.Println("You do not have permissions to delete the user.")
	case http.StatusNotFound:
		fmt.Println("User with given ID does not exist.")
	default:
		fmt.Println("Internal server error, try again later.")
	}
}

func init() {
	UserCmd.AddCommand(UserRegisterCmd)
	UserCmd.AddCommand(UserDeleteCmd)
	UserCmd.AddCommand(UserAuthCmd)
	UserCmd.AddCommand(UserActivateCmd)
}
