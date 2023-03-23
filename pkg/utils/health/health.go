package health

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var alive bool = true
var last_words string = ""

func health(w http.ResponseWriter, r *http.Request) {
	status := ""
	if alive {
		status = "alive"
	} else {
		status = "dead"
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, fmt.Sprintf("{\"status\": \"%s\", \"last_words\": \"%s\"}", status, strings.ReplaceAll(last_words, "\"", "\\\"")))
}

func Dead(cause string) {
	alive = false
	last_words = fmt.Sprintf("%s%s |||", last_words, cause)
}

func RunServer() {
	http.HandleFunc("/health", health)
	server := &http.Server{
		Addr:         ":9999",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	go server.ListenAndServe()
}
