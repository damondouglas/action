package dialogflow

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/damondouglas/assert"
)

func TestSanitize(t *testing.T) {
	a := assert.New(t)
	data, err := ioutil.ReadFile(".mock/request.json")
	a.HandleError(err)
	transData, err := sanitize(data)
	a.HandleError(err)
	unsanitized, _ := regexp.Compile("param[1,2]\\.original")
	sanitized, _ := regexp.Compile("param[1,2]original")
	a.Equals(unsanitized.Match(data), true)
	a.Equals(sanitized.Match(data), false)
	a.Equals(unsanitized.Match(transData), false)
	a.Equals(sanitized.Match(transData), true)
}
func TestRequest(t *testing.T) {
	var req *Request
	a := assert.New(t)
	data, err := ioutil.ReadFile(".mock/request.json")
	a.HandleError(err)
	r := bytes.NewBuffer(data)
	req, err = Encode(r)
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
