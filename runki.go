package main

type LookupResult struct {
	Transcript string
	Meanings   []string
}

func main() {
	ya := YandexProvider{
		ApiKey:   YA_API_KEY,
		LangPair: "en-ru",
	}

	ya.Lookup("eloquent")

	return
}
