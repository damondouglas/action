package dialogflow

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// ActionHandler is the action handler method signature.
type ActionHandler func(w http.ResponseWriter, r *http.Request, actionRequest *Request)

var (
	handlers map[string]ActionHandler
)

func init() {
	handlers = map[string]ActionHandler{}
}

// HandleAction registers ActionHandler with name.
// name is derived from Request.QueryResult.Action when parsed by Dispatch.
// name is the same as the intent name in dialogflow.
func HandleAction(name string, handler ActionHandler) {
	handlers[name] = handler
}

// Dispatch is the main handler provided to http.HandlerFunc("...", Dispatch)
func Dispatch(w http.ResponseWriter, r *http.Request) {
	var req *Request

	req, err := Encode(r.Body)
	if err != nil {
		log.Println(errors.WithStack(err))
	}

	action := req.QueryResult.Action
	if action == "" {
		log.Println("action is blank")
	}
	if f, ok := handlers[action]; ok {
		f(w, r, req)
	} else {
		log.Println("handler not found for action " + action)
	}
}
