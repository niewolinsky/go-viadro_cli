package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type DocumentList struct {
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
	Document bytes.Buffer `json:"document"`
	Metadata string       `json:"metadata"`
}

// * COMMANDS * //
var DocumentCmd = &cobra.Command{
	Use:   "document",
	Short: "Manage documents",
	Long:  "Manage documents",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available subcommands: list, get, delete, grab, merge, toggle, upload")
	},
}
var DocumentGetCmd = &cobra.Command{
	Use:     "get <document_id>",
	Example: "viadro document get 1 --details",
	Short:   "Get details about a document",
	Long:    "Get details about a document",
	Run:     cmdGetDocument,
	Args:    cobra.ExactArgs(1),
}
var DocumentDeleteCmd = &cobra.Command{
	Use:     "delete <document_id>",
	Example: "viadro document delete 1",
	Short:   "Delete document with given id",
	Long:    "Delete document with given id",
	Run:     cmdDeleteDocument,
	Args:    cobra.ExactArgs(1),
}
var DocumentListCmd = &cobra.Command{
	Use:     "list",
	Example: "viadro document list",
	Short:   "List all visible (public) documents",
	Long:    "List all visible (public) documents",
	Run:     cmdGetAllDocuments,
}
var DocumentToggleCmd = &cobra.Command{
	Use:     "toggle <document_id>",
	Example: "viadro document toggle 1",
	Short:   "Toggle document visibility",
	Long:    "Toggle document visibility",
	Run:     cmdToggleDocument,
	Args:    cobra.ExactArgs(1),
}
var DocumentUploadCmd = &cobra.Command{
	Use:     "upload <filepath> <tags>",
	Example: "viadro document upload sample.pdf tag1,tag2",
	Short:   "Upload document (.txt, .docx, .pdf, .rtf, .md) to the cloud",
	Long:    "Upload document (.txt, .docx, .pdf, .rtf, .md) to the cloud",
	Run:     cmdUploadDocument,
	Args:    cobra.ExactArgs(2),
}
var DocumentMergeCmd = &cobra.Command{
	Use:     "merge <title> <tags> <filepath1> <filepath2> ...",
	Example: "viadro document merge mergedpdf work,school sample1.pdf sample2.pdf",
	Short:   "Merge many documents into one and upload to the cloud",
	Long:    "Merge many documents into one and upload to the cloud",
	Run:     cmdMergeDocuments,
	Args:    cobra.MinimumNArgs(3),
}
var DocumentGrabCmd = &cobra.Command{
	Use:     "grab <title> <tags> <link>",
	Example: "viadro document grab grabbedpdf work,school example.com/sample.pdf",
	Short:   "Pull document from the web and host it on Viadro service",
	Long:    "Pull document from the web and host it on Viadro service",
	Run:     cmdGrabDocument,
	Args:    cobra.ExactArgs(3),
}

// * RUN * //
func cmdGetDocument(cmd *cobra.Command, args []string) {
	document_id, err := strconv.Atoi(args[0])
	if err != nil {
		Logger.Fatal("invalid document id")
	}

	urlPart := viper.GetString("endpoint")
	url := urlPart + fmt.Sprintf("/document/%d", document_id)

	showDetails, err := cmd.Flags().GetBool("details")
	if err != nil {
		Logger.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
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
			Logger.Fatal("app error")
		}

		Logger.Info("requested document:")
		if showDetails {
			fmt.Printf("ID: %d | TITLE: %s | LINK: %s | TAGS: %v | UPLOADED: %v | HIDDEN: %v \n", respStruct.Document.DocumentID, respStruct.Document.Title, respStruct.Document.URLS3, respStruct.Document.Tags, respStruct.Document.UploadedAt, respStruct.Document.IsHidden)
		} else {
			fmt.Printf("TITLE: %s | LINK: %s | TAGS: %v \n", respStruct.Document.Title, respStruct.Document.URLS3, respStruct.Document.Tags)
		}
	case http.StatusUnauthorized:
		Logger.Fatal("you do not have permissions to view the document")
	case http.StatusNotFound:
		Logger.Fatal("document with given id does not exist")
	default:
		Logger.Fatal("internal server error, try again later")
	}
}

func cmdDeleteDocument(cmd *cobra.Command, args []string) {
	msg := Delete(args)

	//! not always info, fix
	Logger.Info(msg)
}

func cmdGetAllDocuments(cmd *cobra.Command, args []string) {
	owner, err := cmd.Flags().GetString("owner")
	if err != nil {
		Logger.Fatal(err)
	}

	respStruct := List(args, owner)
	Logger.Info("list of documents: ")
	for _, document := range respStruct.Documents {
		fmt.Printf("ID: %d | TITLE: %s | LINK: %s | TAGS: %v | UPLOADED: %v \n", document.DocumentID, document.Title, document.Link, document.Tags, document.CreatedAt)
	}
}

func cmdToggleDocument(cmd *cobra.Command, args []string) {
	msg := Toggle(args)

	//! not always info, fix
	Logger.Info(msg)
}

func cmdUploadDocument(cmd *cobra.Command, args []string) {
	isHidden, err := cmd.Flags().GetBool("hidden")
	if err != nil {
		Logger.Fatal("app error")
	}

	r, err := os.ReadFile(args[0])
	if err != nil {
		Logger.Fatal("app error")
	}

	buffer := bytes.NewBuffer(r)

	filename := ""
	if len(args[0]) < 4 || args[0][len(args[0])-4:] != ".pdf" {
		filename = args[0]+".pdf"
	} else {
		filename = args[0]
	}

	UploadDocument(buffer, args[1], isHidden, filename)
}

func cmdMergeDocuments(cmd *cobra.Command, args []string) {
	buffer := bytes.NewBuffer([]byte{})
	isHidden, err := cmd.Flags().GetBool("hidden")
	if err != nil {
		Logger.Fatal("app error")
	}

	err = pdfcpu.Merge(args[2], args[3:], buffer, nil)
	if err != nil {
		Logger.Fatal("app error")
	}

	UploadDocument(buffer, args[1], isHidden, fmt.Sprintf("%s%s", args[0], ".pdf"))
}

func cmdGrabDocument(cmd *cobra.Command, args []string) {
	buffer := bytes.NewBuffer([]byte{})

	isHidden, err := cmd.Flags().GetBool("hidden")
	if err != nil {
		Logger.Fatal("app error")
	}

	url := args[2]
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		Logger.Fatal("app error")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		Logger.Fatal("requested link is unreachable or invalid")
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		_, err = io.Copy(buffer, res.Body)
		if err != nil {
			Logger.Fatal(err)
		}

		UploadDocument(buffer, args[1], isHidden, fmt.Sprintf("%s%s", args[0], ".pdf"))

	default:
		Logger.Fatal("app error")
	}
}

// * FUNCTIONS * //
func List(args []string, owner string) DocumentList {
	URL := viper.GetString("endpoint") + "/documents"

	switch owner {
	case "all":
		URL += "/?owner=all"
	case "me":
		URL += "/?owner=me"
	case "exclude":
		URL += "/?owner=-me"
	case "-me":
		URL += "/?owner=-me"
	}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		Logger.Fatal("app error")
	}

	//! BROKEN ON TYPESCRIPT,
	// if owner != "all" {
		bearer := "Bearer " + viper.GetString("tkn")
		req.Header.Add("Authorization", bearer)
	// }

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		Logger.Fatal("service unavailable, try again later")
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		respStruct := DocumentList{}

		err = json.NewDecoder(res.Body).Decode(&respStruct)
		if err != nil {
			Logger.Fatal("app error")
		}

		return respStruct
	default:
		Logger.Fatal("app error")
	}

	return DocumentList{}
}

func Toggle(args []string) string {
	document_id, err := strconv.Atoi(args[0])
	if err != nil {
		Logger.Fatal("app error")
	}

	urlPart := viper.GetString("endpoint")
	url := urlPart + fmt.Sprintf("/document/%d", document_id)

	req, err := http.NewRequest(http.MethodPatch, url, nil)
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
	case http.StatusOK:
		return fmt.Sprintf("successfully toggled visibility of document with id: %d", document_id)
	case http.StatusUnauthorized:
		return "you do not have permissions to view the document"
	case http.StatusNotFound:
		return "document with given id does not exist"
	default:
		return "internal server error, try again later"
	}
}

func Delete(args []string) string {
	document_id, err := strconv.Atoi(args[0])
	if err != nil {
		Logger.Fatal("app error")
	}

	urlPart := viper.GetString("endpoint")
	url := urlPart + fmt.Sprintf("/document/%d", document_id)

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
	case http.StatusOK:
		return fmt.Sprintf("successfully deleted document with id: %d", document_id)
	case http.StatusUnauthorized:
		return "you do not have permissions to delete the document"
	case http.StatusNotFound:
		return "document with given id does not exist"
	default:
		return "internal server error, try again later"
	}
}

func UploadDocument(file *bytes.Buffer, tagsX string, isHidden bool, title string) {
	client := &http.Client{Timeout: 10 * time.Second}
	URL := viper.GetString("endpoint") + "/document"

	//* FIX
	tagz := strings.Split(tagsX, ",")
	tags := fmt.Sprintf("%#v", tagz)
	tags = tags[8:]
	tags = strings.Replace(tags, "{", "[", -1)
	tags = strings.Replace(tags, "}", "]", -1)

	metadata := fmt.Sprintf(`{"is_hidden": %v, "tags": %s}`, isHidden, tags)

	input := Multipart{
		Document: *file,
		Metadata: metadata,
	}

	Upload(client, URL, input, title)
}

func Upload(client *http.Client, url string, input Multipart, title string) (err error) {
	b := bytes.Buffer{}
	w := multipart.NewWriter(&b)

	mtype := mimetype.Detect(input.Document.Bytes())
	if mtype.String() != "application/pdf" || mtype.Extension() != ".pdf" {
		Logger.Fatal(errors.New("invalid file format, for now Viadro only accepts PDF files"))
	}

	wr, err := w.CreateFormFile("document", title)
	if err != nil {
		Logger.Fatal("app error")
	}

	_, err = io.Copy(wr, &input.Document)
	if err != nil {
		Logger.Fatal("app error")
	}

	wr, err = w.CreateFormField("metadata")
	if err != nil {
		Logger.Fatal("app error")
	}

	_, err = io.Copy(wr, strings.NewReader(input.Metadata))
	if err != nil {
		Logger.Fatal("app error")
	}
	w.Close()

	bearer := "Bearer " + viper.GetString("tkn")
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		Logger.Fatal("app error")
	}

	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		Logger.Fatal("server not responding, try again later")
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusCreated {
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
			Logger.Fatal("app error")
		}

		Logger.Info("document was uploaded |", "id", respStruct.Document.DocumentID, "link", respStruct.Document.URLS3, "hidden", respStruct.Document.IsHidden)
	} else if res.StatusCode == http.StatusUnauthorized {
		Logger.Fatal("invalid or expired token, use auth command to grab a new token")
	} else {
		Logger.Fatal("app error")
	}

	return
}

// * INIT * //
func init() {
	DocumentCmd.AddCommand(DocumentGetCmd)
	DocumentGetCmd.PersistentFlags().Bool("details", false, "See file details? Default: hidden")

	DocumentCmd.AddCommand(DocumentListCmd)
	DocumentListCmd.PersistentFlags().StringP("owner", "o", "all", "Files from which owners to show: all, me, exclude/-me")

	DocumentCmd.AddCommand(DocumentToggleCmd)

	DocumentCmd.AddCommand(DocumentDeleteCmd)

	DocumentCmd.AddCommand(DocumentUploadCmd)
	DocumentUploadCmd.PersistentFlags().Bool("hidden", false, "Should the file be hidden? Default: visible")

	DocumentCmd.AddCommand(DocumentMergeCmd)
	DocumentMergeCmd.PersistentFlags().Bool("hidden", false, "Should the file be hidden? Default: visible")

	DocumentCmd.AddCommand(DocumentGrabCmd)
	DocumentGrabCmd.PersistentFlags().Bool("hidden", false, "Should the file be hidden? Default: visible")
}
