package dialogflow

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"path"
	"regexp"
)

// Request data from dialogflow fulfillment webhook.
type Request struct {
	QueryResult struct {
		Action                   string      `json:"action"`
		QueryText                string      `json:"queryText"`
		Parameters               *Parameters `json:"parameters"`
		AllRequiredParamsPresent bool        `json:"allRequiredParamsPresent"`
		OutputContexts           []struct {
			Name       string      `json:"name"`
			Parameters *Parameters `json:"parameters"`
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
	transData := re.ReplaceAll(data, []byte("original"))
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
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&values); err != nil {
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
	var num json.Number
	var ok bool

	if num, ok = value.(json.Number); !ok {
		strValue := value.(string)
		param := new(stringParam)
		param.value = &strValue
		return param
	}

	intValue, _ := num.Int64()
	floatValue, _ := num.Float64()
	if floatValue-float64(intValue) > 0 {
		param := new(floatParam)
		param.value = &floatValue
		return param
	}

	param := new(intParam)
	param.value = &intValue
	return param
}

type stringParam struct {
	value *string
}

func (param *stringParam) Get() interface{} {
	return *param.value
}

type intParam struct {
	value *int64
}

func (param *intParam) Get() interface{} {
	return int(*param.value)
}

type floatParam struct {
	value *float64
}

func (param *floatParam) Get() interface{} {
	return *param.value
}

// Value is the parameter interface.
type Value interface {
	Get() interface{}
}
