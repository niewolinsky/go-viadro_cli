package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// userCmd represents the user command
var AdminUserListCmd = &cobra.Command{
	Use:   "list",
	Short: "Manage users and documents as administrator",
	Long:  ``,
	Run:   listAllAdmin,
}

func listAllAdmin(cli *cobra.Command, args []string) {
	URL := "http://localhost:4000/v1/admin/users"
	bearer := "Bearer " + viper.GetString("tkn")

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal("Can't form request")
	}
	req.Header.Add("Authorization", bearer)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Service unavailable, try again later.")
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		respStruct := struct {
			Users []struct {
				ID        int       `json:"id"`
				CreatedAt time.Time `json:"created_at"`
				Username  string    `json:"username"`
				Email     string    `json:"email"`
				Activated bool      `json:"activated"`
				IsAdmin   bool      `json:"is_admin"`
			} `json:"users"`
		}{}

		err = json.NewDecoder(res.Body).Decode(&respStruct)
		if err != nil {
			fmt.Println(err)
		}

		for _, user := range respStruct.Users {
			fmt.Println(user)
		}
	} else if res.StatusCode == http.StatusUnauthorized {
		fmt.Println("Invalid or expired token, use auth command to grab a new token.")
	} else {
		fmt.Println("internal server error, try again later", res.StatusCode)
	}
}

func init() {
	AdminUserCmd.AddCommand(AdminUserListCmd)
}
