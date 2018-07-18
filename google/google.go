package google

import (
	"errors"
	"unicode/utf8"
)

// Google is the response required in the webhook fulfillment.
type Google struct {
	ExpectUserResponse bool          `json:"expectUserResponse"`
	RichResponse       *RichResponse `json:"richResponse"`
}

// RichResponse can include audio, text, cards, suggestions and structured data Items.
type RichResponse struct {
	Items       []Item        `json:"items"`
	Suggestions []*Suggestion `json:"suggestions"`
}

// Item of the response.
type Item struct {
	SimpleResponse *SimpleResponse `json:"simpleResponse,omitempty"`
}

// SimpleResponse contains speech or text to show the user.
type SimpleResponse struct {
	TextToSpeech string `json:"textToSpeech"`
	DisplayText  string `json:"displayText"`
}

// Suggestion chip that the user can tap to quickly post a reply to the conversation.
type Suggestion struct {
	Title string `json:"title"`
}

// Response returns Google fulfillment data.
func Response(expectUserResponse bool, items []Item, suggestions []string) *Google {
	rr := &RichResponse{
		Items: items,
	}
	if len(suggestions) > 0 {
		s := make([]*Suggestion, len(suggestions))
		for i, v := range suggestions {
			s[i] = &Suggestion{
				Title: v,
			}
		}
		rr.Suggestions = s
	}
	return &Google{
		ExpectUserResponse: expectUserResponse,
		RichResponse:       rr,
	}
}

// Simple provides Item with SimpleResponse.
func Simple(show *DisplayText, say string) Item {
	return Item{
		SimpleResponse: &SimpleResponse{
			TextToSpeech: say,
			DisplayText:  show.value,
		},
	}
}

// DisplayText restricts string to 640 characters based on google assistant limitation.
type DisplayText struct {
	value string
}

// NewDisplayText validates value to 640 character limit.
func NewDisplayText(value string) (*DisplayText, error) {
	if utf8.RuneCountInString(value) > 640 {
		return nil, errors.New("display text exceeds character limit")
	}
	return &DisplayText{
		value: value,
	}, nil
}
