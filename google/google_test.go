package google

import "log"

func ExampleSimple() {
	text := "Hello there"
	show, err := NewDisplayText(text)
	if err != nil {
		log.Panicln(err)
	}
	say := "hi"
	item := Simple(show, say)

	log.Println(item)
}
