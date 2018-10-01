package main

import (
	"fmt"
	"github.com/BlockClusterApp/daemon/src/helpers"
	"github.com/BlockClusterApp/daemon/src/tools"
	"net/http"

	"github.com/BlockClusterApp/daemon/src/config"
	"github.com/gorilla/mux"
)

var log *helpers.Logger

func main() {
	log = helpers.GetLogger()
	router := newRouter()
	tools.StartScheduler()

	log.Fatal(http.ListenAndServe(":3005", router))
}

func newRouter() *mux.Router {
	router := mux.NewRouter()

	// log("Config %s", config.GetKubeConfig)

	router.HandleFunc("/ping", handlePing).Methods("GET")
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
	fmt.Fprintf(w, "%s", config.GetKubeConfig())
}
