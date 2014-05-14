package main

import (
	"fmt"
)

type LookupResult struct {
	Transcript string
	Meanings   []string
}

func main() {
	//ya := YandexProvider{
	//    ApiKey:   YA_API_KEY,
	//    LangPair: "en-ru",
	//}

	anki, _ := AnkiWebLogin("login", "pass")
	fmt.Println(anki)

	//fmt.Println(anki.Search("eloquent"))
	fmt.Println(anki.Add("test", "eloquent", "a"))

	return
}
