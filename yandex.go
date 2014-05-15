package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const YA_API_URI = "https://dictionary.yandex.net" +
	"/api/v1/dicservice.json/lookup?key=%s&lang=%s&text=%s"
const YA_API_KEY = "dict.1.1.20140512T122957Z.549af1de13649236." +
	"74bbc11e0fa7625166dd95f21b1ff17838df2c03"

type YandexProvider struct {
	key  string
	lang string
}

func NewYandexProvider(lang, key string) *YandexProvider {
	if key == "" {
		key = YA_API_KEY
	}
	return &YandexProvider{key, lang}
}

func (y YandexProvider) Lookup(text string) (*LookupResult, error) {
	resp, err := http.Get(fmt.Sprintf(YA_API_URI, y.key, y.lang, text))
	if err != nil {
		return nil, err
	}

	if resp.Status != "200 OK" {
		return nil, errors.New(
			"expected HTTP status 200, received %s" + resp.Status)
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
			}
		}
	}{}

	json.Unmarshal(body, &dictResult)

	lookupResult := LookupResult{
		Transcript: dictResult.Def[0].Ts,
		Meanings:   make([]string, len(dictResult.Def)),
	}

	for i, d := range dictResult.Def {
		lookupResult.Meanings[i] = d.Tr[0].Text
	}

	return &lookupResult, nil
}
