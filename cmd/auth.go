/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "A brief description of your command",
	Long:  ``,
	Run:   auth,
}

func auth(cmd *cobra.Command, args []string) {
	URL := "http://localhost:4000/v1/users/authenticate"

	input := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    "user1@gmail.com",
		Password: "haslo456",
	}

	jsonified, _ := json.MarshalIndent(input, "", "\t")
	reader := bytes.NewReader(jsonified)

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest(http.MethodPut, URL, reader)
	if err != nil {
		fmt.Println("error tutaj")
	}

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("asd")
	}
	defer response.Body.Close()

	respStruct := struct {
		AuthenticationToken struct {
			Token  string    `json:"token"`
			Expiry time.Time `json:"expiry"`
		} `json:"authentication_token"`
	}{}

	err = json.NewDecoder(response.Body).Decode(&respStruct)
	if err != nil {
		fmt.Println(err)
	}

	os.Setenv(respStruct.AuthenticationToken.Token, "viadro_auth_token")
	fmt.Println("Sucessfully authenticated with token: ", respStruct.AuthenticationToken.Token)
}

func init() {
	rootCmd.AddCommand(authCmd)
}
