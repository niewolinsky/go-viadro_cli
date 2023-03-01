package cli

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"viadro_cli/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var UserActivateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Activate user account",
	Long:  ``,
	Run:   userActivate,
	Args:  cobra.ExactArgs(1),
}

func userActivate(cli *cobra.Command, args []string) {
	input := struct {
		TokenPlaintext string `json:"token"`
	}{
		TokenPlaintext: args[0],
	}

	req := utils.PrepareRequest(input, viper.GetString("endpoint")+"/users/activate", http.MethodPut)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Service unavailable, try again later.")
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		fmt.Println("User successfully activated, use auth command to login.")
	} else if res.StatusCode == http.StatusBadRequest {
		fmt.Println("malformed json request", res.StatusCode)
	} else if res.StatusCode == http.StatusUnprocessableEntity {
		fmt.Println("invalid or expired token", res.StatusCode)
	} else {
		fmt.Println("internal server error, try again later", res.StatusCode)
	}
}

func init() {
	UserCmd.AddCommand(UserActivateCmd)
}
