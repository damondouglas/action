package dialogflow

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"regexp"
	"strconv"
)

// Request data from dialogflow fulfillment webhook.
type Request struct {
	QueryResult struct {
		QueryText                string      `json:"queryText"`
		Action                   string      `json:"action"`
		Parameters               *Parameters `json:"parameters"`
		AllRequiredParamsPresent bool        `json:"allRequiredParamsPresent"`
		OutputContexts           []struct {
			Name string `json:"name"`
		} `json:"outputContexts"`
	} `json:"queryResult"`
}

// Contexts extracts context names from dialogflow request.
func (r *Request) Contexts() []string {
	ctxs := make([]string, len(r.QueryResult.OutputContexts))
	for i, c := range r.QueryResult.OutputContexts {
		ctxs[i] = path.Base(c.Name)
	}
	return ctxs
}

func sanitize(data []byte) ([]byte, error) {
	// dialogflow sends json key as <paramname>.original and encoder ignores it.
	re, err := regexp.Compile("\\.original")
	if err != nil {
		return nil, err
	}
	transData := re.ReplaceAll(data, []byte("org"))
	return transData, nil
}

// Encode Request from Reader.
func Encode(r io.Reader) (*Request, error) {
	var req *Request
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	transData, err := sanitize(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(transData, &req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// Parameters parses intent parameters from dialogflow request.
type Parameters struct {
	values map[string]Value
}

// UnmarshalJSON is an implementation of json.UnmarshalJSON
func (p *Parameters) UnmarshalJSON(data []byte) error {
	var values map[string]interface{}
	err := json.Unmarshal(data, &values)
	if err != nil {
		return err
	}
	p.values = make(map[string]Value, len(values))
	for k, v := range values {
		p.values[k] = cast(v)
	}
	return nil
}

// Get parameter value from dialogflow request by key name.
// The key matches the parameter name from dialogflow intent.
func (p *Parameters) Get(key string) interface{} {
	if val, ok := p.values[key]; ok {
		return val.Get()
	}
	return nil
}

func cast(value interface{}) Value {
	var val Value

	strValue := fmt.Sprint(value)
	defaultParam := new(stringParam)
	defaultParam.value = &strValue
	val = defaultParam
	intValue, err := strconv.Atoi(strValue)
	if err == nil {
		param := new(intParam)
		param.value = &intValue
		val = param
	}

	return val
}

type stringParam struct {
	value *string
}

func (param *stringParam) Get() interface{} {
	return *param.value
}

type intParam struct {
	value *int
}

func (param *intParam) Get() interface{} {
	return *param.value
}

// Value is the parameter interface.
type Value interface {
	Get() interface{}
}
