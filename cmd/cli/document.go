package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	Document *os.File `json:"document"`
	Metadata string   `json:"metadata"`
}

// * COMMANDS * //
var DocumentCmd = &cobra.Command{
	Use:   "document",
	Short: "Manage documents",
	Long:  "Manage documents",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Available subcommands: list, get, grab, merge, toggle, upload")
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
	Use:     "merge <filepath1> <filepath2> ...",
	Example: "viadro document merge sample1.pdf sample2.pdf",
	Short:   "Merge many documents into one and upload to the cloud",
	Long:    "Merge many documents into one and upload to the cloud",
	Run:     cmdDocumentsMerge,
	Args:    cobra.MinimumNArgs(2),
}
var DocumentGrabCmd = &cobra.Command{
	Use:     "grab <link>",
	Example: "viadro document grab example.com/sample.pdf",
	Short:   "Pull document from the web and host it on Viadro service",
	Long:    "Pull document from the web and host it on Viadro service",
	Run:     cmdGrabDocument,
	Args:    cobra.ExactArgs(1),
}

// * RUN * //
func cmdGetDocument(cmd *cobra.Command, args []string) {
	document_id, err := strconv.Atoi(args[0])
	if err != nil {
		logger.Fatal("invalid document id")
	}

	url := fmt.Sprintf(`http://localhost:4000/v1/document/%d`, document_id)
	bearer := "Bearer " + viper.GetString("tkn")

	showDetails, err := cmd.Flags().GetBool("details")
	if err != nil {
		logger.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
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
			logger.Fatal("app error")
		}

		logger.Info("requested document:")
		if showDetails {
			fmt.Printf("ID: %d | TITLE: %s | LINK: %s | TAGS: %v | UPLOADED: %v | HIDDEN: %v \n", respStruct.Document.DocumentID, respStruct.Document.Title, respStruct.Document.URLS3, respStruct.Document.Tags, respStruct.Document.UploadedAt, respStruct.Document.IsHidden)
		} else {
			fmt.Printf("TITLE: %s | LINK: %s | TAGS: %v \n", respStruct.Document.Title, respStruct.Document.URLS3, respStruct.Document.Tags)
		}
	case http.StatusUnauthorized:
		logger.Fatal("you do not have permissions to view the document")
	case http.StatusNotFound:
		logger.Fatal("document with given id does not exist")
	default:
		logger.Fatal("internal server error, try again later")
	}
}

func cmdGetAllDocuments(cmd *cobra.Command, args []string) {
	owner, err := cmd.Flags().GetString("owner")
	if err != nil {
		logger.Fatal(err)
	}

	respStruct := ListTesting(args, owner)
	logger.Info("list of documents: ")
	for _, document := range respStruct.Documents {
		fmt.Printf("ID: %d | TITLE: %s | LINK: %s | TAGS: %v | UPLOADED: %v \n", document.DocumentID, document.Title, document.Link, document.Tags, document.CreatedAt)
	}
}

func cmdToggleDocument(cmd *cobra.Command, args []string) {
	msg := Toggle(args)

	//! not always info, fix
	logger.Info(msg)
}

func cmdUploadDocument(cmd *cobra.Command, args []string) {
	isHidden, err := cmd.Flags().GetBool("hidden")
	if err != nil {
		logger.Fatal("app error")
	}

	UploadDocument(args[0], args[1], isHidden)
}

func cmdDocumentsMerge(cmd *cobra.Command, args []string) {
	mergeTempFile, err := os.Create("merge-temp.pdf")
	if err != nil {
		logger.Fatal("app error")
	}
	defer mergeTempFile.Close()

	err = pdfcpu.Merge(args[0], args[1:], mergeTempFile, nil)
	if err != nil {
		logger.Fatal("app error")
	}

	UploadDocument("merge-temp.pdf", "merged,test", false)
}

func cmdGrabDocument(cmd *cobra.Command, args []string) {
	// Create the file
	out, err := os.Create("grab-temp.pdf")
	if err != nil {
		logger.Fatal("app error")
	}
	defer out.Close()

	url := args[0]
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.Fatal("app error")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		logger.Fatal("requested link is unreachable or invalid")
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		_, err = io.Copy(out, res.Body)
		if err != nil {
			logger.Fatal(err)
		}

		UploadDocument("grab-temp.pdf", "grab,net", false)

	default:
		logger.Fatal("app error")
	}
}

// * FUNCTIONS * //
func ListTesting(args []string, owner string) DocumentList {
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
		logger.Fatal("app error")
	}

	if owner != "all" {
		bearer := "Bearer " + viper.GetString("tkn")
		req.Header.Add("Authorization", bearer)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		logger.Fatal("service unavailable, try again later")
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		respStruct := DocumentList{}

		err = json.NewDecoder(res.Body).Decode(&respStruct)
		if err != nil {
			logger.Fatal("app error")
		}

		return respStruct
	default:
		logger.Fatal("app error")
	}

	return DocumentList{}
}

func Toggle(args []string) string {
	document_id, err := strconv.Atoi(args[0])
	if err != nil {
		logger.Fatal("app error")
	}

	url := fmt.Sprintf(`http://localhost:4000/v1/document/%d`, document_id)

	req, err := http.NewRequest(http.MethodPatch, url, nil)
	if err != nil {
		logger.Fatal("app error")
	}

	bearer := "Bearer " + viper.GetString("tkn")
	req.Header.Add("Authorization", bearer)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		logger.Fatal("service unavailable, try again later")
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

func MustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}

func UploadDocument(filepath string, tagsX string, isHidden bool) {
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
		Document: MustOpen(filepath),
		Metadata: metadata,
	}

	Upload(client, URL, input)
}

func Upload(client *http.Client, url string, input Multipart) (err error) {
	b := bytes.Buffer{}
	w := multipart.NewWriter(&b)
	wr, err := w.CreateFormFile("document", input.Document.Name())
	if err != nil {
		logger.Fatal(err)
	}

	_, err = io.Copy(wr, input.Document)
	if err != nil {
		logger.Fatal(err)
	}

	wr, err = w.CreateFormField("metadata")
	if err != nil {
		logger.Fatal(err)
	}

	_, err = io.Copy(wr, strings.NewReader(input.Metadata))
	if err != nil {
		logger.Fatal(err)
	}
	w.Close()

	bearer := "Bearer " + viper.GetString("tkn")
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		logger.Fatal(err)
	}
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", w.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
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
			logger.Fatal("app error")
		}

		logger.Info("document was uploaded |", "id", respStruct.Document.DocumentID, "link", respStruct.Document.URLS3, "hidden", respStruct.Document.IsHidden)
	} else if res.StatusCode == http.StatusUnauthorized {
		logger.Fatal("invalid or expired token, use auth command to grab a new token")
	} else {
		logger.Fatal("app error")
	}

	return
}

// * INIT * //
func init() {
	DocumentCmd.AddCommand(DocumentGetCmd)
	DocumentGetCmd.PersistentFlags().Bool("details", false, "See file details? Default: hidden")

	DocumentCmd.AddCommand(DocumentListCmd)
	DocumentListCmd.PersistentFlags().StringP("owner", "o", "all", "Files from which owners to show: all, me, exclude")

	DocumentCmd.AddCommand(DocumentToggleCmd)

	DocumentCmd.AddCommand(DocumentUploadCmd)
	DocumentUploadCmd.PersistentFlags().Bool("hidden", false, "Should the file be hidden? Default: visible")

	DocumentCmd.AddCommand(DocumentMergeCmd)

	DocumentCmd.AddCommand(DocumentGrabCmd)
}
