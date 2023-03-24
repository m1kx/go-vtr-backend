package health

import (
	"fmt"
	"net/http"
	"strings"
)

var alive bool = true
var last_words string = ""

func Health(w http.ResponseWriter, r *http.Request) {
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
