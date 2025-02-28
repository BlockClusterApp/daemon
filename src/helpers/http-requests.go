package helpers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/getsentry/raven-go"
	"io/ioutil"
	"net/http"
)

func MakeHTTPRequest(url string, method string, payload string) (string, error){
	//GetLogger().Printf("Sending body %s", payload)
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	bc := GetBlockclusterInstance()
	//log.Println("Auth", bc.AuthToken)
	if bc.AuthToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(bc.AuthToken))))
	} else {
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("fetch-token"))))
	}


	var client = &http.Client{}


	resp, err := client.Do(req)

	if err != nil {
		raven.CaptureError(err, map[string]string{
			"licenceKey": bc.Licence.Key,
			"url": url,
			"method": method,
			"payload": payload,
		})
		GetLogger().Printf("Error making request: %s", err.Error())
		return "", err // Can't find cert file
	}
	defer resp.Body.Close()

	if resp.StatusCode > 401 {
		raven.CaptureError(err, map[string]string{
			"licenceKey": bc.Licence.Key,
			"url": url,
			"method": method,
			"payload": payload,
		})
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		GetLogger().Printf("Request to %s returned %d %s", url, resp.StatusCode, bodyBytes)
		resp.Body.Close()
		return "", errors.New(fmt.Sprintf("Request to %s returned %d", url, resp.StatusCode))
	}

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)

		if err2 != nil {
			GetLogger().Printf("Error reading body for %s %s", url, err2.Error())
			return "",err2
		}

		return bodyString, nil
	} else if resp.StatusCode == http.StatusUnauthorized {
		Blockcluster.AuthToken = ""
		Blockcluster.Reauthorize()
	}

	return "",errors.New(fmt.Sprintf("Unhandled status code for %s | %s", url, resp.Status))

}