package main

import (
	"fmt"
	"github.com/BlockClusterApp/daemon/src/helpers"
	"github.com/BlockClusterApp/daemon/src/tools"
	"github.com/getsentry/raven-go"
	"net/http"
	"os"

	"github.com/BlockClusterApp/daemon/src/config"
	"github.com/gorilla/mux"
)

var log *helpers.Logger

func main() {
	raven.SetDSN("https://3fb09492cf20449aae350ac935dcd26d:8b2fbf9455f34ce08fd51ba7e8042919@sentry.io/1302256")

	log = helpers.GetLogger()
	router := newRouter()
	tools.StartScheduler()

	log.Fatal(http.ListenAndServe(":3005", router))
}

func newRouter() *mux.Router {
	router := mux.NewRouter()

	// log("Config %s", config.GetKubeConfig)

	router.HandleFunc("/ping", handlePing).Methods("GET")
	router.HandleFunc("/healthz", handlePing).Methods("GET")
	router.HandleFunc("/config", handleConfig).Methods("GET")

	return router
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	//log.Println("Handle /ping")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", "Pong")
}

func handleConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("Handle /config")

	// In kubernetes cluster, always verify the license
	if os.Getenv("GO_ENV") == "development" && os.Getenv("KUBERNETES_SERVICE_PORT_HTTPS") == "" {
		fmt.Fprintf(w, "%s", config.GetKubeConfig())
		return
	}
	var bc = helpers.GetBlockclusterInstance()
	if bc.Valid == false {
		fmt.Fprintf(w, "%s", "{\"error\": \"Licence Invalid\", \"errorCode\": 404}")
		return
	}
	fmt.Fprintf(w, "%s", config.GetKubeConfig())
}
