[![godoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/damondouglas/action)

# About

Package action is a golang client for dialog flow requests and responses through Google Actions.

# Usage

```golang
dialogflow.HandleAction("foo", func(w http.ResponseWriter, r *http.Request, req *dialogflow.Request) {
    param1 := req.QueryResult.Parameters.Get("param1")
    param2 := req.QueryResult.Parameters.Get("param2")
    text := fmt.Sprintf("param1: %v\nparam2: %v", param1, param2)
    show, err := google.NewDisplayText(text)
    if err != nil {
        log.Println(errors.WithStack(err))
    }
    resp := dialogflow.Google(true, []string{"foo"}, google.Simple(show, text))
    json.NewEncoder(w).Encode(resp)
})

http.HandleFunc("/action", dialogflow.Dispatch)
```
