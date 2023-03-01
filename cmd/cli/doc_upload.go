package cli

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var UploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload document to the cloud",
	Long:  ``,
	Run:   documentUpload,
	Args:  cobra.ExactArgs(2),
}

type Multipart struct {
	Document *os.File `json:"document"`
	Metadata string   `json:"metadata"`
}

func documentUpload(cli *cobra.Command, args []string) {
	client := &http.Client{Timeout: 10 * time.Second}
	URL := viper.GetString("endpoint") + "/document"

	is_hidden, err := cli.Flags().GetBool("hidden")
	if err != nil {
		log.Fatal(err)
	}

	//* UGLY AF, FIX
	tagz := strings.Split(args[1], ",")
	tags := fmt.Sprintf("%#v", tagz)
	tags = tags[8:]
	tags = strings.Replace(tags, "{", "[", -1)
	tags = strings.Replace(tags, "}", "]", -1)

	metadata := fmt.Sprintf(`{"is_hidden": %v, "tags": %s}`, is_hidden, tags)

	input := Multipart{
		Document: mustOpen(args[0]),
		Metadata: metadata,
	}

	Upload(client, URL, input)

	// values := map[string]io.Reader{
	// 	// "document": mustOpen(args[0]),
	// 	"metadata": strings.NewReader(fmt.Sprintf(`{"is_hidden": %v, "tags": %v}`, is_hidden, tags)),
	// }

	// fmt.Println(values["metadata"])

	// err = Upload(client, URL, values)
	// if err != nil {
	// 	log.Fatal(err)
	// }
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

func init() {
	UploadCmd.PersistentFlags().Bool("hidden", false, "Should the file be hidden? Default: visible")
}
