package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// * COMMANDS * //
var AdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Manage users and documents as administrator",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available subcommands: user, document")
	},
}
var AdminUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users and documents as administrator",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available subcommands: delete, list, grant")
	},
}
var AdminUserListCmd = &cobra.Command{
	Use:   "list",
	Short: "Manage users and documents as administrator",
	Run:   cmdGetAllUsers,
}
var AdminUserGrantCmd = &cobra.Command{
	Use:   "grant",
	Short: "Manage users and documents as administrator",
	Run:   cmdGrantAdmin,
	Args:  cobra.ExactArgs(1),
}
var AdminDocumentCmd = &cobra.Command{
	Use:   "document",
	Short: "Manage users and documents as administrator",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available subcommands: list")
	},
}
var AdminDocumentListCmd = &cobra.Command{
	Use:   "list",
	Short: "Manage users and documents as administrator",
	Run:   cmdGetAllDocumentsAdmin,
}

// * RUN * //
func cmdGetAllUsers(cmd *cobra.Command, args []string) {
	URL := "http://localhost:4000/v1/admin/users"
	bearer := "Bearer " + viper.GetString("tkn")

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal("can't form request", "error", err)
	}
	req.Header.Add("Authorization", bearer)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("service unavailable, try again later")
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
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
	case http.StatusUnauthorized:
		fmt.Println("invalid or expired token, use auth command to grab a new token")
	default:
		fmt.Println("internal server error, try again later")
	}
}

func cmdGrantAdmin(cmd *cobra.Command, args []string) {
	user_id, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf(`http://localhost:4000/v1/admin/user/%d`, user_id)

	req, err := http.NewRequest(http.MethodPatch, url, nil)
	if err != nil {
		log.Fatal("service unavailable, try again later")
	}

	bearer := "Bearer " + viper.GetString("tkn")
	req.Header.Add("Authorization", bearer)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("service unavailable, try again later")
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		respStruct := struct {
			User struct {
				UserID    int       `json:"user_id"`
				CreatedAt time.Time `json:"created_at"`
				Username  string    `json:"username"`
				Email     string    `json:"email"`
				Activated bool      `json:"activated"`
				IsAdmin   bool      `json:"is_admin"`
			} `json:"user"`
		}{}

		err = json.NewDecoder(res.Body).Decode(&respStruct)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Successfully toggled admin privileges of user with ID: %d", respStruct.User.UserID)
	case http.StatusUnauthorized:
		fmt.Println("you do not have permissions to toggle admin privileges")
	case http.StatusNotFound:
		fmt.Println("user with given ID does not exist")
	default:
		fmt.Println("internal server error, try again later")
	}
}

func cmdGetAllDocumentsAdmin(cmd *cobra.Command, args []string) {
	URL := "http://localhost:4000/v1/admin/documents"

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal("can't form request")
	}

	bearer := "Bearer " + viper.GetString("tkn")
	req.Header.Add("Authorization", bearer)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("service unavailable, try again later")
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		respStruct := struct {
			Documents []struct {
				DocumentID int       `json:"document_id"`
				Title      string    `json:"title"`
				Link       string    `json:"link"`
				Tags       []string  `json:"tags"`
				CreatedAt  time.Time `json:"created_at"`
			} `json:"documents"`
			Metadata struct {
				CurrentPage  int `json:"current_page"`
				PageSize     int `json:"page_size"`
				FirstPage    int `json:"first_page"`
				LastPage     int `json:"last_page"`
				TotalRecords int `json:"total_records"`
			} `json:"metadata"`
		}{}

		err = json.NewDecoder(res.Body).Decode(&respStruct)
		if err != nil {
			fmt.Println(err)
		}

		for _, document := range respStruct.Documents {
			fmt.Println(document)
		}

	case http.StatusUnauthorized:
		fmt.Println("invalid or expired token, use auth command to grab a new token")
	case http.StatusForbidden:
		fmt.Println("you dont have required role or privilege to complete this action")
	default:
		fmt.Println("internal server error, try again later")
	}
}

// * INIT * //
func init() {
	AdminCmd.AddCommand(AdminUserCmd)
	AdminUserCmd.AddCommand(AdminUserListCmd)
	AdminUserCmd.AddCommand(AdminUserGrantCmd)

	AdminCmd.AddCommand(AdminDocumentCmd)
	AdminDocumentCmd.AddCommand(AdminDocumentListCmd)
}
