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
	Long:  "Manage users and documents as administrator",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available subcommands: user, document")
	},
}
var AdminUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage users as administrator",
	Long:  "Manage users as administrator",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available subcommands: delete, list, grant")
	},
}
var AdminUserListCmd = &cobra.Command{
	Use:     "list",
	Example: "viadro admin user list",
	Short:   "List all registered users",
	Long:    "List all registered users",
	Run:     cmdGetAllUsers,
}
var AdminUserGrantCmd = &cobra.Command{
	Use:     "grant <user_id>",
	Example: "viadro admin user grant 1",
	Short:   "Grant user admin privileges",
	Long:    "Grant user admin privileges",
	Run:     cmdGrantAdmin,
	Args:    cobra.ExactArgs(1),
}
var AdminDocumentCmd = &cobra.Command{
	Use:   "document",
	Short: "Manage documents as administrator",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available subcommands: list")
	},
}
var AdminDocumentListCmd = &cobra.Command{
	Use:     "list",
	Example: "viadro admin document list",
	Short:   "List all documents regardless of document visibility",
	Long:    "List all documents regardless of document visibility",
	Run:     cmdGetAllDocumentsAdmin,
}

// * RUN * //
func cmdGetAllUsers(cmd *cobra.Command, args []string) {
	URL := "http://localhost:4000/v1/admin/users"
	bearer := "Bearer " + viper.GetString("tkn")

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		logger.Fatal("app error")
	}
	req.Header.Add("Authorization", bearer)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		logger.Fatal("service unavailable, try again later")
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		respStruct := struct {
			Users []struct {
				UserId    int       `json:"user_id"`
				CreatedAt time.Time `json:"created_at"`
				Username  string    `json:"username"`
				Email     string    `json:"email"`
				Activated bool      `json:"activated"`
				IsAdmin   bool      `json:"is_admin"`
			} `json:"users"`
		}{}

		err = json.NewDecoder(res.Body).Decode(&respStruct)
		if err != nil {
			logger.Fatal("app error")
		}

		logger.Info("list of users: ")
		for _, user := range respStruct.Users {
			fmt.Printf("ID: %d | CREATED: %v | USERNAME: %s | EMAIL: %s | ACTIVATED: %v | ADMIN: %v \n", user.UserId, user.CreatedAt, user.Username, user.Email, user.Activated, user.IsAdmin)
		}
	case http.StatusUnauthorized:
		logger.Fatal("invalid or expired token, use auth command to grab a new token")
	case http.StatusForbidden:
		logger.Fatal("you do not have required privileges to perform this action")
	default:
		logger.Fatal("internal server error, try again later")
	}
}

func cmdGrantAdmin(cmd *cobra.Command, args []string) {
	user_id, err := strconv.Atoi(args[0])
	if err != nil {
		logger.Fatal("app error")
	}

	url := fmt.Sprintf(`http://localhost:4000/v1/admin/user/%d`, user_id)

	req, err := http.NewRequest(http.MethodPatch, url, nil)
	if err != nil {
		logger.Fatal("app error")
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
			logger.Fatal("app error")
		}

		logger.Info("Successfully toggled admin privileges of user with id: %d", respStruct.User.UserID)
	case http.StatusUnauthorized:
		logger.Fatal("you do not have permissions to toggle admin privileges")
	case http.StatusNotFound:
		logger.Fatal("user with given id does not exist")
	default:
		logger.Fatal("internal server error, try again later")
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
		respStruct := DocumentList{}

		err = json.NewDecoder(res.Body).Decode(&respStruct)
		if err != nil {
			logger.Fatal("app error")
		}

		logger.Info("list of documents: ")
		for _, document := range respStruct.Documents {
			fmt.Printf("ID: %d | TITLE: %s | LINK: %s | TAGS: %v | UPLOADED: %v \n", document.DocumentID, document.Title, document.Link, document.Tags, document.CreatedAt)
		}

	case http.StatusUnauthorized:
		logger.Fatal("invalid or expired token, use auth command to grab a new token")
	case http.StatusForbidden:
		logger.Fatal("you do not have required privileges to perform this action")
	default:
		logger.Fatal("internal server error, try again later")
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
