package dialogflow

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/damondouglas/action/google"
	"github.com/damondouglas/assert"
	"github.com/pkg/errors"
)

var (
	mockAction = Action{
		Name:    "foo",
		Handler: mockActionHandler,
		RequiredParams: []string{
			"param1",
			"param2",
		},
	}
)

func TestDispatch(t *testing.T) {
	a := assert.New(t)

	HandleAction(mockAction)
	_, ok := handlers["foo"]
	a.Equals(ok, true)

	f, _ := os.Open(".mock/request.json")
	req := httptest.NewRequest("GET", "/action", f)
	w := httptest.NewRecorder()

	Dispatch(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	a.Equals(string(body), "foo frog 3")
}

func mockActionHandler(w http.ResponseWriter, r *http.Request, req *Request) {
	fmt.Fprint(w, req.QueryResult.QueryText)
}

func Example() {
	handler := func(w http.ResponseWriter, r *http.Request, req *Request) {
		param1 := req.QueryResult.Parameters.Get("param1")
		param2 := req.QueryResult.Parameters.Get("param2")
		text := fmt.Sprintf("param1: %v\nparam2: %v", param1, param2)
		show, err := google.NewDisplayText(text)
		if err != nil {
			log.Println(errors.WithStack(err))
		}
		resp := Google(true, []string{"foo"}, google.Simple(show, text))
		json.NewEncoder(w).Encode(resp)

	}
	fooAction := Action{
		Name:    "foo",
		Handler: handler,
		RequiredParams: []string{
			"param1",
			"param2",
		},
	}
	HandleAction(fooAction)

	http.HandleFunc("/action", Dispatch)
}

func TestHandleAction(t *testing.T) {
	var req *Request
	a := assert.New(t)

	HandleAction(mockAction)
	_, ok := handlers["foo"]
	a.Equals(ok, true)

	f, _ := os.Open(".mock/request.json")
	w := httptest.NewRecorder()

	err := json.NewDecoder(f).Decode(&req)
	a.HandleError(err)

	mockActionHandler(w, nil, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	a.Equals(string(body), "foo frog 3")
}
