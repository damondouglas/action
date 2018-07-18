package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/damondouglas/action/dialogflow"
	"github.com/damondouglas/action/google"
	"github.com/pkg/errors"
	"google.golang.org/appengine"
)

func main() {
	http.HandleFunc("/action", dialogflow.Dispatch)
	dialogflow.HandleAction("foo", foo)
	dialogflow.HandleAction("input.welcome", welcome)
	appengine.Main()
}

func action(w http.ResponseWriter, r *http.Request) {
}

func foo(w http.ResponseWriter, r *http.Request, req *dialogflow.Request) {
	param1 := req.QueryResult.Parameters.Get("param1")
	param2 := req.QueryResult.Parameters.Get("param2")
	text := fmt.Sprintf("param1: %v\nparam2: %v", param1, param2)
	show, err := google.NewDisplayText(text)
	if err != nil {
		log.Println(errors.WithStack(err))
	}
	resp := dialogflow.Google(true, []string{"foo"}, google.Simple(show, text))
	json.NewEncoder(w).Encode(resp)
}

func welcome(w http.ResponseWriter, r *http.Request, req *dialogflow.Request) {
	text := "Welcome to showcase. What would you like to see?"
	show, err := google.NewDisplayText(text)
	if err != nil {
		log.Println(errors.WithStack(err))
	}
	resp := dialogflow.Google(true, []string{"foo"}, google.Simple(show, text))
	json.NewEncoder(w).Encode(resp)
}
