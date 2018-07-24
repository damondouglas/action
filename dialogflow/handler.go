package dialogflow

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
)

const originalKey = "original"

// ActionHandler is the action handler method signature.
type ActionHandler func(w http.ResponseWriter, r *http.Request, actionRequest *Request)

var (
	handlers map[string]Action
)

func init() {
	handlers = map[string]Action{}
}

// HandleAction registers ActionHandler with name.
// name is derived from Request.QueryResult.Action when parsed by Dispatch.
// name is the same as the intent name in dialogflow.
func HandleAction(action Action) {
	handlers[action.Name] = action
}

// Dispatch is the main handler provided to http.HandlerFunc("...", Dispatch)
func Dispatch(w http.ResponseWriter, r *http.Request) {
	var req *Request
	var action Action
	var ok bool
	var err error

	req, err = Encode(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(errors.WithStack(err))
	}

	name := req.QueryResult.Intent.DisplayName
	if name == "" {
		log.Println("action is blank")
	}
	if action, ok = handlers[name]; !ok {
		log.Println("handler not found for action " + name)
	}
	if ok = action.requestIsValid(req); !ok {
		log.Println("all required parameters are not present in request. maybe dialogflow action is not configured with output context")
	}
	f := action.Handler
	f(w, r, req)
}

// Action configures dialogflow action response.
type Action struct {
	Name           string
	Handler        ActionHandler
	RequiredParams []string
}

func (a *Action) requestIsValid(req *Request) bool {
	params := req.QueryResult.Parameters
	originalParams := req.QueryResult.Parameters
	for _, paramName := range a.RequiredParams {
		if value := params.Get(paramName); value == "" {
			return false
		}
		if value := originalParams.Get(paramName + originalKey); value == "" {
			return false
		}
	}
	return true
}
