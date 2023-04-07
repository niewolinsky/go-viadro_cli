package cli

import (
	"encoding/json"
	"fmt"
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
	Long:  "Manage user's credentials",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available subcommands: register, activate, auth|login, delete|remove|deactivate")
	},
}
var UserRegisterCmd = &cobra.Command{
	Use:     "register <username> <email> <password>",
	Example: "viadro user register user1 user1@mail.com pass",
	Long:    "Register a new user for the Viadro service",
	Short:   "Register a new user for the Viadro service",
	Run:     cmdUserRegister,
	Args:    cobra.ExactArgs(3),
}
var UserActivateCmd = &cobra.Command{
	Use:     "activate <token>",
	Example: "viadro user activate TH2AP6G6IWWZA7MQRLO6F4C2PI",
	Short:   "Activate user account",
	Long:    "Activate user account",
	Run:     cmdUserActivate,
	Args:    cobra.ExactArgs(1),
}
var UserAuthCmd = &cobra.Command{
	Use:     "auth <email> <password>",
	Aliases: []string{"login"},
	Example: "viadro user auth user1 pass",
	Short:   "Authenticate (login) an existing user",
	Long:    "Authenticate (login) an existing user",
	Run:     cmdUserAuth,
	Args:    cobra.ExactArgs(2),
}
var UserDeleteCmd = &cobra.Command{
	Use:     "delete <user_id>",
	Aliases: []string{"remove, deactivate"},
	Example: "viadro user delete 1",
	Short:   "Delete (deactivate) a user from Viadro service",
	Long:    "Delete (deactivate) a user from Viadro service",
	Run:     cmdUserDelete,
	Args:    cobra.ExactArgs(1),
}
var UserLogoutCmd = &cobra.Command{
	Use:     "logout",
	Example: "viadro user logout",
	Short:   "Logout from the service",
	Long:    "Logout from the service",
	Run:     cmdUserLogout,
}

// * RUN * //
func cmdUserRegister(cmd *cobra.Command, args []string) {
	input := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Username: args[0],
		Email:    args[1],
		Password: args[2],
	}

	req := utils.PrepareRequest(input, viper.GetString("endpoint")+"/user", http.MethodPost)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		Logger.Fatal("service unavailable, try again later")
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusAccepted:
		respStruct := struct {
			User struct {
				UserId    int       `json:"user_id"`
				CreatedAt time.Time `json:"created_at"`
				Username  string    `json:"username"`
				Email     string    `json:"email"`
				Activated bool      `json:"activated"`
			} `json:"user"`
		}{}

		err = json.NewDecoder(res.Body).Decode(&respStruct)
		if err != nil {
			Logger.Fatal("app error")
		}

		Logger.Info("created new user, check your email for account activation", "user_id", respStruct.User.UserId)
	case http.StatusBadRequest:
		Logger.Fatal("malformed json request")
	case http.StatusUnprocessableEntity:
		Logger.Fatal("account with given email exists")
	default:
		Logger.Fatal("app error")
	}
}

func cmdUserActivate(cmd *cobra.Command, args []string) {
	input := struct {
		TokenPlaintext string `json:"token"`
	}{
		TokenPlaintext: args[0],
	}

	req := utils.PrepareRequest(input, viper.GetString("endpoint")+"/user/activate", http.MethodPut)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		Logger.Fatal("service unavailable, try again later")
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		Logger.Info("user successfully activated, use auth command to login")
	case http.StatusBadRequest:
		Logger.Fatal("malformed json request")
	case http.StatusUnprocessableEntity:
		Logger.Fatal("invalid or expired token")
	default:
		Logger.Fatal("app error")
	}
}

func cmdUserAuth(cmd *cobra.Command, args []string) {
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
		Logger.Fatal("service unavailable, try again later")
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
			Logger.Fatal("app error")
		}

		Logger.Info("successfully logged in |", "token", respStruct.AuthenticationToken.Token)
		Logger.Info("token has been stored in config, for the next 24 hours requests will automatically use it")
	case http.StatusBadRequest:
		Logger.Fatal("malformed json request")
	case http.StatusUnauthorized:
		Logger.Fatal("wrong email or password")
	default:
		Logger.Fatal("app error")
	}
}

func cmdUserDelete(cmd *cobra.Command, args []string) {
	user_id, err := strconv.Atoi(args[0])
	if err != nil {
		Logger.Fatal("invalid user id")
	}

	urlPart := viper.GetString("endpoint")
	url := urlPart + fmt.Sprintf("/user/%d", user_id)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		Logger.Fatal("app error")
	}

	bearer := "Bearer " + viper.GetString("tkn")
	req.Header.Add("Authorization", bearer)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		Logger.Fatal("service unavailable, try again later")
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusNoContent:
		fmt.Println("user and its files successfully deleted")
	case http.StatusUnauthorized:
		Logger.Fatal("you do not have the permissions to delete the user")
	case http.StatusNotFound:
		Logger.Fatal("user with given id does not exists")
	default:
		Logger.Fatal("app error")
	}
}

func cmdUserLogout(cmd *cobra.Command, args []string) {
	token := viper.GetString("tkn")

	if token == "" {
		Logger.Error("user is not logged in")
	} else {
		viper.Set("tkn", "")
		err := viper.WriteConfig()
		if err != nil {
			Logger.Fatal("app error")
		}
		Logger.Info("user logged out successfully")
	}
}

// * INIT * //
func init() {
	UserCmd.AddCommand(UserRegisterCmd)

	UserCmd.AddCommand(UserDeleteCmd)

	UserCmd.AddCommand(UserAuthCmd)

	UserCmd.AddCommand(UserActivateCmd)

	UserCmd.AddCommand(UserLogoutCmd)
}
