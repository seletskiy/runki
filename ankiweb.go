package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
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

var ANKI_RE_MID = regexp.MustCompile(`mid":\s*"?(\d+)`)
var ANKI_RE_ITEM = regexp.MustCompile(`(?s:mitem3.*?<td>([^/]+))`)

func NewAnkiAccount() *AnkiAccount {
	return &AnkiAccount{}
}

func (a *AnkiAccount) Load(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	storedData := struct {
		Cookies []*http.Cookie
		Mid     string
	}{
		make([]*http.Cookie, 0),
		"",
	}
	err = json.Unmarshal(data, &storedData)
	if err != nil {
		return err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}

	ankiurl, _ := url.Parse(ANKI_HOST)
	jar.SetCookies(ankiurl, storedData.Cookies)

	a.http = &http.Client{Jar: jar}
	a.mid = storedData.Mid

	return nil
}

func (a *AnkiAccount) Save(filename string) error {
	ankiurl, _ := url.Parse(ANKI_HOST)
	data, err := json.Marshal(struct {
		Cookies []*http.Cookie
		Mid     string
	}{
		a.http.Jar.Cookies(ankiurl),
		a.mid,
	})
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(filename), 0700)
	if err != nil {
		fmt.Println(err)
	}

	ioutil.WriteFile(filename, data, 0600)

	return nil
}

func (a *AnkiAccount) WebLogin(login, password string) error {
	ankiurl, _ := url.Parse(ANKI_HOST)

	jar, _ := cookiejar.New(nil)
	a.http = &http.Client{Jar: jar}

	resp, err := a.http.Get(ANKI_HOST)
	_, err = a.http.PostForm(ANKI_LOGIN_URI,
		url.Values{
			"username": {login},
			"password": {password},
		})

	if err != nil {
		return err
	}

	if len(jar.Cookies(ankiurl)) == 0 {
		return errors.New("failed to login to anki web")
	}

	resp, _ = a.http.Get(ANKI_EDITOR_URI)
	body, _ := ioutil.ReadAll(resp.Body)
	midMatch := ANKI_RE_MID.FindStringSubmatch(string(body))

	if len(midMatch) < 2 {
		return errors.New("failed to get mid authentification value from anki web")
	}

	a.mid = midMatch[1]

	return nil
}

func (a *AnkiAccount) Search(searchText string) (bool, error) {
	resp, err := a.http.PostForm(ANKI_SEARCH_URI,
		url.Values{
			"keyword":   {searchText},
			"submitted": {"1"},
		})

	if err != nil {
		return false, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	items := ANKI_RE_ITEM.FindAllStringSubmatch(string(body), -1)
	for _, m := range items {
		if strings.TrimSpace(m[1]) == searchText {
			return true, nil
		}
	}

	return false, nil
}

func (a *AnkiAccount) Add(deck, text, translation string) error {
	data, err := json.Marshal([]interface{}{
		[]string{text, translation},
		"", // tag is unused
	})

	urlValues := url.Values{
		"data": {string(data)},
		"mid":  {a.mid},
		"deck": {deck},
	}

	resp, err := a.http.PostForm(ANKI_ADD_URI, urlValues)

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
