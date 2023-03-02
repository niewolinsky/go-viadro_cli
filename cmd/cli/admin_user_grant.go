package cli

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// userCmd represents the user command
var AdminUserGrantCmd = &cobra.Command{
	Use:   "grant",
	Short: "Manage users and documents as administrator",
	Long:  ``,
	Run:   grantAdmin,
	Args:  cobra.ExactArgs(1),
}

func grantAdmin(cmd *cobra.Command, args []string) {
	user_id, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf(`http://localhost:4000/v1/admin/user/%d`, user_id)
	bearer := "Bearer " + viper.GetString("tkn")

	req, err := http.NewRequest(http.MethodPatch, url, nil)
	if err != nil {
		log.Fatal("Service unavailable, try again later.")
	}

	req.Header.Add("Authorization", bearer)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Service unavailable, try again later.")
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)

	//TODO: FINISH
}

func init() {
	AdminUserCmd.AddCommand(AdminUserGrantCmd)
}
