package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	UnlimitedSynonyms = 0
)

// flags=4 is for morphological search
const (
	YandexApiUrlFormat = "https://dictionary.yandex.net" +
		"/api/v1/dicservice.json/lookup?key=%s&lang=%s&text=%s&flags=4"
	YandexApiKey = "dict.1.1.20140512T122957Z.549af1de13649236." +
		"74bbc11e0fa7625166dd95f21b1ff17838df2c03"
)

type YandexProvider struct {
	key           string
	lang          string
	limitSynonyms int
}

func NewYandexProvider(lang, key string, limitSynonyms int) *YandexProvider {
	if key == "" {
		key = YandexApiKey
	}

	return &YandexProvider{key, lang, limitSynonyms}
}

func (y YandexProvider) Lookup(text string) (*LookupResult, error) {
	url := fmt.Sprintf(YandexApiUrlFormat, y.key, y.lang, url.QueryEscape(text))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(
			fmt.Sprintf("expected HTTP status 200, received %s (%s)",
				resp.Status, url))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	dictResult := struct {
		Def []struct {
			Pos string
			Ts  string
			Tr  []struct {
				Text string
				Mean []struct {
					Text string
				}
			}
		}
	}{}

	err = json.Unmarshal(body, &dictResult)
	if err != nil {
		return nil, err
	}

	if len(dictResult.Def) == 0 {
		return nil, nil
	}

	lookupResult := LookupResult{
		Transcript: dictResult.Def[0].Ts,
	}

	for _, d := range dictResult.Def {
		for j, tr := range d.Tr {
			if y.limitSynonyms > 0 && j >= y.limitSynonyms {
				break
			}
			references := []string{}
			for _, ref := range tr.Mean {
				references = append(references, ref.Text)
			}

			lookupResult.Meanings = append(lookupResult.Meanings, Meaning{
				Translation: tr.Text,
				References:  references,
			})
		}
	}

	return &lookupResult, nil
}
