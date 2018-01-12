package controllers

import (
	"fmt"
	"net/http"
	"os"

	cfg "github.com/mccraymt/ms-black-history/config"
)

// HandleStatusRequest reutrns basic heartbeat/status check data
func HandleStatusRequest(w http.ResponseWriter, r *http.Request) {
	host, _ := os.Hostname()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ms-black-history version %v\n", cfg.Config.Version)
	fmt.Fprintf(w, "%v %v\n", cfg.Config.Environment, host)
}

// SayError ...
func SayError(w http.ResponseWriter, r *http.Request) {
	//e.ReturnError(w, e.GenericError, "something went wrong")
	host, _ := os.Hostname()
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "ms-black-history version %v\n", cfg.Config.Version)
	fmt.Fprintf(w, "%v %v\n", cfg.Config.Environment, host)
}
