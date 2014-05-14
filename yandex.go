package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type YaDictResult struct {
	Def []YaDef `json:"def"`
}

type YaDef struct {
	Pos         string `json:"pos"`
	Transcript  string `json:"ts"`
	Translation []YaTr `json:"tr"`
}

type YaTr struct {
	Text string `json:"text"`
}

const YA_API_URI = "https://dictionary.yandex.net/api/v1/dicservice.json/lookup?key=%s&lang=%s&text=%s"
const YA_API_KEY = "dict.1.1.20140512T122957Z.549af1de13649236.74bbc11e0fa7625166dd95f21b1ff17838df2c03"

type YandexProvider struct {
	ApiKey   string
	LangPair string
}

func (y YandexProvider) Lookup(text string) (*LookupResult, error) {
	resp, err := http.Get(fmt.Sprintf(YA_API_URI, y.ApiKey, y.LangPair, text))
	if err != nil {
		return nil, err
	}

	if resp.Status != "200 OK" {
		return nil, errors.New("expected HTTP status 200, received %s" + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	dictResult := YaDictResult{}

	json.Unmarshal(body, &dictResult)

	fmt.Println(dictResult)

	return nil, nil
}
