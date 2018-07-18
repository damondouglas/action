package dialogflow

import (
	"github.com/damondouglas/action/google"
)

// Response is required for dialogflow fulfillment.
type Response struct {
	Payload *Payload `json:"payload"`
}

// Payload is required to route dialogflow fulfillment.
type Payload struct {
	Google *google.Google `json:"google,omitempty"`
}

// Google returns dialogflow fulfillment response for use with google assistant actions.
func Google(expectUserResponse bool, suggestions []string, items ...google.Item) *Response {
	payload := &Payload{
		Google: google.Response(expectUserResponse, items, suggestions),
	}
	return &Response{
		Payload: payload,
	}
}
