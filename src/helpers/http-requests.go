package helpers

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"os"
)

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

func MakeRequest(method string,path string, payload io.Reader) (string, error){
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
			log.Printf("Cert not found: %s", err.Error())
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