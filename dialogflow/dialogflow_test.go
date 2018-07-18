package dialogflow

import (
	"log"

	"github.com/damondouglas/action/google"
)

func ExampleGoogle() {
	text := "Hello there"
	show, err := google.NewDisplayText(text)
	if err != nil {
		log.Panicln(err)
	}
	say := "hi"
	item := google.Simple(show, say)

	expectResponse := true
	suggestions := []string{
		"some suggestion",
	}

	resp := Google(expectResponse, suggestions, item)
	log.Println(resp)
}
