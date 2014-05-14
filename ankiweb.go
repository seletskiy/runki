package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

var _ = fmt.Print

type AnkiAccount struct {
	http *http.Client
	mid  string
}

const ANKI_HOST = "https://ankiweb.net"
const ANKI_LOGIN_URI = ANKI_HOST + "/account/login"
const ANKI_SEARCH_URI = ANKI_HOST + "/search/"
const ANKI_ADD_URI = ANKI_HOST + "/edit/save"
const ANKI_EDITOR_URI = ANKI_HOST + "/edit/"

var ANKI_RE_MID = regexp.MustCompile(`mid":\s*"(\d+)"`)
var ANKI_RE_ITEM = regexp.MustCompile(`(?s:mitem3.*?<td>([^/]+))`)

func AnkiWebLogin(login, password string) (*AnkiAccount, error) {
	ankiurl, _ := url.Parse(ANKI_HOST)

	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar}

	_, err := client.Get(ANKI_HOST)
	_, err = client.PostForm(ANKI_LOGIN_URI,
		url.Values{
			"username": {login},
			"password": {password},
		})

	if err != nil {
		return nil, err
	}

	if len(jar.Cookies(ankiurl)) == 0 {
		return nil, errors.New("failed to login to anki web")
	}

	resp, _ := client.Get(ANKI_EDITOR_URI)
	body, _ := ioutil.ReadAll(resp.Body)
	midMatch := ANKI_RE_MID.FindStringSubmatch(string(body))

	if len(midMatch) < 2 {
		return nil, errors.New("failed to get mid from anki web")
	}

	return &AnkiAccount{&client, midMatch[1]}, nil
}

func (a AnkiAccount) Search(search_text string) (bool, error) {
	resp, err := a.http.PostForm(ANKI_SEARCH_URI,
		url.Values{
			"keyword":   {search_text},
			"submitted": {"1"},
		})

	if err != nil {
		return false, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	items := ANKI_RE_ITEM.FindAllStringSubmatch(string(body), -1)
	for _, m := range items {
		if strings.TrimSpace(m[1]) == search_text {
			return true, nil
		}
	}

	return false, nil
}

func (a AnkiAccount) Add(deck, text, translation string) error {
	data, err := json.Marshal([]interface{}{
		[]string{text, translation},
		"", // tag is unused
	})

	resp, err := a.http.PostForm(ANKI_ADD_URI,
		url.Values{
			"data": {string(data)},
			"mid":  {a.mid},
			"deck": {deck},
		})

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	bodyStr := string(body)

	if bodyStr != "1" {
		return errors.New("unexpected answer from anki web while add: " +
			bodyStr)
	}

	return nil
}
