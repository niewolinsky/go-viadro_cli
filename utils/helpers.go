package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func PrepareRequest(input interface{}, url string, method string) *http.Request {
	jsonified, _ := json.MarshalIndent(input, "", "\t")
	reader := bytes.NewReader(jsonified)

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		log.Fatal("can't form request")
	}

	req.Header.Set("Content-Type", "application/json")

	return req
}
