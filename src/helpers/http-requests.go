package helpers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func MakeHTTPRequest(url string, method string, payload string) (string, error){
	log := GetLogger()
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	bc := GetBlockclusterInstance()
	if bc.AuthToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(bc.AuthToken))))
	} else {
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("fetch-token"))))
	}


	var client = &http.Client{}


	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Error making request: %s", err.Error())
		return "", err // Can't find cert file
	}
	defer resp.Body.Close()

	if resp.StatusCode > 400 {
		log.Print("Request to %s returned %d", url, resp.StatusCode)
		resp.Body.Close()
		return "", errors.New(fmt.Sprintf("Request to %s returned %d", url, resp.StatusCode))
	}

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)

		if err2 != nil {
			log.Printf("Error reading body for %s", url, err2.Error())
			return "",err2
		}

		return bodyString, nil
	}

	return "",errors.New(fmt.Sprintf("Unhandled status code for %s | %s", url, resp.Status))

}