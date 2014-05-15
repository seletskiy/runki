package main

import (
	"fmt"
)

type LookupResult struct {
	Transcript string
	Meanings   []string
}

func main() {
	ya := NewYandexProvider("en-ru", "")
	anki := NewAnkiAccount("~/.runki/creds")

	anki.WebLogin("s.seletskiy@gmail.com", "zmxncbv13")
	fmt.Println(anki)

	fmt.Println(anki.Add("test", "eloquent", "a"))

	return
}
