package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/damondouglas/action/dialogflow"
	"github.com/damondouglas/action/google"
	"github.com/pkg/errors"
	"google.golang.org/appengine"
)

func main() {
	http.HandleFunc("/action", action)
	appengine.Main()
}

func action(w http.ResponseWriter, r *http.Request) {
	show, err := google.NewDisplayText("new\nline")
	if err != nil {
		log.Println(errors.WithStack(err))
	}
	resp := dialogflow.Google(true, []string{"Here's a suggestion"}, google.Simple(show, "I'm bold and I'm italic."))
	json.NewEncoder(w).Encode(resp)
}
