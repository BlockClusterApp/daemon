package helpers

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/getsentry/raven-go"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"os"
)

type ExternalKubeRequest struct {
	URL  string
	Auth struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	}
	Method  string
	Payload string
}

func isInKubernetes() bool {
	serviceHost := os.Getenv("KUBERNETES_SERVICE_HOST")
	if len(serviceHost) == 0 {
		return false
	}
	return true
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func getURL(path string) string {

	serviceHost := os.Getenv("KUBERNETES_SERVICE_HOST")
	servicePort := os.Getenv("KUBERNETES_SERVICE_PORT")

	var url string

	if !isInKubernetes() {
		kubeApiServer := os.Getenv("KUBE_API_SERVER_URL")
		if len(kubeApiServer) == 0 {
			log.Fatal("KUBE_API_SERVER_URL and KUBERNETES_SERVICE_HOST both are not present in env")
			return "";
		}
		url = fmt.Sprintf("%s%s", kubeApiServer, path)
	} else {
		url = fmt.Sprintf("https://%s:%s%s", serviceHost, servicePort, path)
	}

	u, err := url2.Parse(url)
	if err != nil {
		panic(err)
	}

	return u.String()
}

func MakeExternalKubeRequest(params ExternalKubeRequest) (string, error) {
	req, err := http.NewRequest(params.Method, params.URL, bytes.NewBuffer([]byte(params.Payload)))

	if err != nil {
		GetLogger().Printf("Error creating external kube request %s %s", params.URL, err.Error())
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicAuth(params.Auth.User, params.Auth.Pass)))

	var client = &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		GetLogger().Printf("Error making external kube request %s %s", params.URL, err.Error())
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)

		resp.Body.Close()
		if err2 != nil {
			GetLogger().Printf("Error reading body for %s %s", params.URL, err2.Error())
			return "", err2
		}

		return bodyString, nil
	}

	resp.Body.Close()
	return "", errors.New(fmt.Sprintf("Unhandled status code for %s | %s", params.URL, resp.Status))
}

func MakeKubeRequest(method string, path string, payload io.Reader) (string, error) {
	var url string;
	url = getURL(path)

	req, err := http.NewRequest(method, url, payload)
	req.Header.Set("Content-Type", "application/json")

	var client = &http.Client{}

	if isInKubernetes() {
		caToken, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
		if err != nil {
			panic(err) // cannot find token file
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", string(caToken)))

		caCertPool := x509.NewCertPool()
		caCert, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
		if err != nil {
			GetLogger().Printf("Cert not found: %s", err.Error())
			return "", err // Can't find cert file
		}

		caCertPool.AppendCertsFromPEM(caCert)

		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: caCertPool,
				},
			},
		}
	} else {
		// Not in kubernetes
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicAuth(os.Getenv("KUBE_API_USER"), os.Getenv("KUBE_API_PASS"))))
	}

	resp, err := client.Do(req)

	if err != nil {
		bc := GetBlockclusterInstance()
		raven.CaptureError(err, map[string]string{
			"licenceKey": bc.Licence.Key,
		})
		GetLogger().Printf("Error making request: %s", err.Error())
		return "", err // Can't find cert file
	}
	defer resp.Body.Close()

	if resp.StatusCode > 400 {
		GetLogger().Printf("Request to %s returned %d", url, resp.StatusCode)
		resp.Body.Close()
		return "", errors.New(fmt.Sprintf("Request to %s returned %d", url, resp.StatusCode))
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		bodyBytes, err2 := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)

		resp.Body.Close()
		if err2 != nil {
			GetLogger().Printf("Error reading body for %s %s", url, err2.Error())
			return "", err2
		}

		return bodyString, nil
	}

	resp.Body.Close()
	return "", errors.New(fmt.Sprintf("Unhandled status code for %s | %s", url, resp.Status))

}
