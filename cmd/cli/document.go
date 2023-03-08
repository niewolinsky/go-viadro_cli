package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RespStruct struct {
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
}

type Multipart struct {
	Document *os.File `json:"document"`
	Metadata string   `json:"metadata"`
}

// userCmd represents the user command
var DocumentCmd = &cobra.Command{
	Use:   "document",
	Short: "Manage documents",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available subcommands: list, get, grab, merge, toggle, upload")
	},
}
var DocumentGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details about a document",
	Run:   cmdGetDocument,
	Args:  cobra.ExactArgs(1),
}
var DocumentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all visible (public) documents",
	Run:   cmdGetAllDocuments,
}
var DocumentToggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggle document visibility",
	Run:   cmdToggleDocument,
	Args:  cobra.ExactArgs(1),
}
var DocumentUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload document to the cloud",
	Run:   cmdUploadDocument,
	Args:  cobra.ExactArgs(2),
}
var DocumentMergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "",
	Run:   cmdDocumentsMerge,
	Args:  cobra.MinimumNArgs(3),
}
var DocumentGrabCmd = &cobra.Command{
	Use:   "grab",
	Short: "Manage documents",
	Run:   cmdGrabDocument,
	Args:  cobra.ExactArgs(1),
}

func cmdGetDocument(cmd *cobra.Command, args []string) {
	document_id, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf(`http://localhost:4000/v1/document/%d`, document_id)
	bearer := "Bearer " + viper.GetString("tkn")

	_, err = cmd.Flags().GetBool("details")
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	switch res.StatusCode {
	case http.StatusOK:
		respStruct := struct {
			Document struct {
				DocumentID int       `json:"document_id"`
				UserID     int       `json:"user_id"`
				URLS3      string    `json:"url_s3"`
				Filetype   string    `json:"filetype"`
				UploadedAt time.Time `json:"uploaded_at"`
				Title      string    `json:"title"`
				Tags       []string  `json:"tags"`
				IsHidden   bool      `json:"is_hidden"`
			} `json:"document"`
		}{}

		err = json.NewDecoder(res.Body).Decode(&respStruct)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Document title: %s", respStruct.Document.Title)
	case http.StatusUnauthorized:
		fmt.Println("you do not have permissions to view the document")
	case http.StatusNotFound:
		fmt.Println("document with given ID does not exist")
	default:
		fmt.Println("internal server error, try again later")
	}
}

func cmdGetAllDocuments(cmd *cobra.Command, args []string) {
	owner, err := cmd.Flags().GetString("owner")
	if err != nil {
		log.Fatal(err)
	}

	respStruct := ListTesting(args, owner)
	for _, document := range respStruct.Documents {
		fmt.Println(document.Title)
	}
}

func ListTesting(args []string, owner string) RespStruct {
	URL := viper.GetString("endpoint") + "/documents"

	switch owner {
	case "all":
		URL += "/?owner=all"
	case "me":
		URL += "/?owner=me"
	case "exclude":
		URL += "/?owner=-me"
	}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal("Can't form request")
	}

	if owner != "all" {
		bearer := "Bearer " + viper.GetString("tkn")
		req.Header.Add("Authorization", bearer)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Service unavailable, try again later.")
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		respStruct := RespStruct{}

		err = json.NewDecoder(res.Body).Decode(&respStruct)
		if err != nil {
			fmt.Println(err)
		}

		return respStruct
	default:
		fmt.Println("internal server error")
	}

	return RespStruct{}
}

func cmdToggleDocument(cmd *cobra.Command, args []string) {
	msg := Toggle(args)
	fmt.Println(msg)
}

func Toggle(args []string) string {
	document_id, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf(`http://localhost:4000/v1/document/%d`, document_id)
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

	switch res.StatusCode {
	case http.StatusOK:
		return fmt.Sprintf("successfully toggled visibility of document with ID: %d", document_id)
	case http.StatusUnauthorized:
		return "you do not have permissions to view the document"
	case http.StatusNotFound:
		return "document with given ID does not exist"
	default:
		return "internal server error, try again later"
	}
}

func cmdUploadDocument(cmd *cobra.Command, args []string) {
	isHidden, err := cmd.Flags().GetBool("hidden")
	if err != nil {
		log.Fatal(err)
	}

	uploadDocument(args[0], args[1], isHidden)
}

func uploadDocument(filepath string, tagsX string, isHidden bool) {
	client := &http.Client{Timeout: 10 * time.Second}
	URL := viper.GetString("endpoint") + "/document"

	//* UGLY AF, FIX
	tagz := strings.Split(tagsX, ",")
	tags := fmt.Sprintf("%#v", tagz)
	tags = tags[8:]
	tags = strings.Replace(tags, "{", "[", -1)
	tags = strings.Replace(tags, "}", "]", -1)

	metadata := fmt.Sprintf(`{"is_hidden": %v, "tags": %s}`, isHidden, tags)

	input := Multipart{
		Document: mustOpen(filepath),
		Metadata: metadata,
	}

	Upload(client, URL, input)
}

func Upload(client *http.Client, url string, input Multipart) (err error) {
	// Prepare a form that you will submit to that URL.
	b := bytes.Buffer{}
	w := multipart.NewWriter(&b)
	wr, err := w.CreateFormFile("document", input.Document.Name())
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(wr, input.Document)
	if err != nil {
		log.Fatal(err)
	}

	wr, err = w.CreateFormField("metadata")
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(wr, strings.NewReader(input.Metadata))
	if err != nil {
		log.Fatal(err)
	}
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	bearer := "Bearer " + viper.GetString("tkn")
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", bearer)
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusCreated {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	} else if res.StatusCode == http.StatusUnauthorized {
		fmt.Println("Invalid or expired token, use auth command to grab a new token.")
	} else {
		fmt.Println("internal server error, try again later", res.StatusCode)
	}

	return
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}

func cmdDocumentsMerge(cmd *cobra.Command, args []string) {
	mergeTempFile, err := os.Create("merge-temp.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer mergeTempFile.Close()

	err = pdfcpu.Merge(args[0], args[1:], mergeTempFile, nil)
	if err != nil {
		log.Fatal(err)
	}

	uploadDocument("merge-temp.pdf", "merged,test", false)
}

func cmdGrabDocument(cmd *cobra.Command, args []string) {
	// Create the file
	out, err := os.Create("grab-temp.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	url := args[0]
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal("Service unavailable, try again later.")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Service unavailable, try again later.")
	}
	defer res.Body.Close()

	// Check server response
	switch res.StatusCode {
	case http.StatusOK:
		_, err = io.Copy(out, res.Body)
		if err != nil {
			log.Fatal(err)
		}

		uploadDocument("grab-temp.pdf", "grab,net", false)
	}
}

func init() {
	DocumentCmd.AddCommand(DocumentGetCmd)
	DocumentGetCmd.PersistentFlags().Bool("details", false, "See file details? Default: hidden")

	DocumentCmd.AddCommand(DocumentListCmd)
	DocumentListCmd.PersistentFlags().StringP("owner", "v", "all", "Files from which owners to show: all, me, exclude")

	DocumentCmd.AddCommand(DocumentToggleCmd)

	DocumentCmd.AddCommand(DocumentUploadCmd)
	DocumentUploadCmd.PersistentFlags().Bool("hidden", false, "Should the file be hidden? Default: visible")

	DocumentCmd.AddCommand(DocumentMergeCmd)

	DocumentCmd.AddCommand(DocumentGrabCmd)
}
