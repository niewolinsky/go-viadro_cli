package cli

// import (
// 	"bytes"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/pdfcpu/pdfcpu/pkg/api"
// 	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
// 	"github.com/spf13/cobra"
// 	"github.com/spf13/viper"
// )

// var MergeCmd = &cobra.Command{
// 	Use:   "merge",
// 	Short: "",
// 	Long:  ``,
// 	Run:   documentsMerge,
// 	Args:  cobra.ExactArgs(2),
// }

// func documentsMerge(cli *cobra.Command, args []string) {
// 	client := &http.Client{Timeout: 10 * time.Second}
// 	URL := viper.GetString("endpoint") + "/document"

// 	//?
// 	b1, err := os.ReadFile("sample1.pdf") // b1 has type []byte
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	rs1 := bytes.NewReader(b1)

// 	b2, err := os.ReadFile("sample2.pdf") // b2 has type []byte
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	rs2 := bytes.NewReader(b2)

// 	rsArr := []io.ReadSeeker{rs1, rs2}

// 	file, err := os.Create("output.pdf")
// 	if err != nil {
// 		return
// 	}
// 	defer file.Close()

// 	err = api.Merge(rsArr, file, pdfcpu.NewDefaultConfiguration())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	//?

// 	//prepare the reader instances to encode
// 	values := map[string]io.Reader{
// 		"document": mustOpen("output.pdf"), // lets assume its this file
// 		"metadata": strings.NewReader(`{"is_hidden": false, "tags": ["work", "school"]}`),
// 	}
// 	err = Upload(client, URL, values)
// 	if err != nil {
// 		panic(err)
// 	}
// }
