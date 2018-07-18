package dialogflow

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/damondouglas/assert"
)

func TestRequest(t *testing.T) {
	var req *Request
	a := assert.New(t)
	data, err := ioutil.ReadFile(".mock/request.json")
	a.HandleError(err)
	err = json.Unmarshal(data, &req)
	a.HandleError(err)
	a.Equals(req.QueryResult.Action, "foo")
	a.Equals(req.QueryResult.AllRequiredParamsPresent, true)
	ctxs := req.Contexts()
	a.Equals(ctxs[0], "output_foo_context")
	param1 := req.QueryResult.Parameters.Get("param1")
	a.Equals(param1, 3)
	param2 := req.QueryResult.Parameters.Get("param2")
	a.Equals(param2, "frog")
}
