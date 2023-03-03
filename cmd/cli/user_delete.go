package cli

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var UserDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete (deactivate) a user from Viadro service",
	Long:  ``,
	Run:   userDelete,
	Args:  cobra.ExactArgs(3),
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

	if res.StatusCode == http.StatusNoContent {
		fmt.Println("User successfully deleted.")
	} else if res.StatusCode == http.StatusUnauthorized {
		fmt.Println("You do not have permissions to delete the user.")
	} else if res.StatusCode == http.StatusNotFound {
		fmt.Println("User with given ID does not exist.")
	} else {
		fmt.Println("Internal server error, try again later.", res.StatusCode)
	}
}

func init() {
	UserCmd.AddCommand(UserDeleteCmd)
}
