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

const (
	AnkiBaseUrl   = "https://ankiweb.net"
	AnkiLoginUrl  = AnkiBaseUrl + "/account/login"
	AnkiSearchUrl = AnkiBaseUrl + "/search/"
	AnkiAddUrl    = AnkiBaseUrl + "/edit/save"
	AnkiEditorUrl = AnkiBaseUrl + "/edit/"
)

var (
	reMid  = regexp.MustCompile(`mid":\s*"?(\d+)`)
	reItem = regexp.MustCompile(`(?s:mitem3.*?<td>([^/]+))`)
)

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

	ankiurl, _ := url.Parse(AnkiBaseUrl)
	jar.SetCookies(ankiurl, storedData.Cookies)

	a.http = &http.Client{Jar: jar}
	a.mid = storedData.Mid

	return nil
}

func (a *AnkiAccount) Save(filename string) error {
	ankiurl, _ := url.Parse(AnkiBaseUrl)
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
	ankiurl, _ := url.Parse(AnkiBaseUrl)

	jar, _ := cookiejar.New(nil)
	a.http = &http.Client{Jar: jar}

	resp, err := a.http.Get(AnkiBaseUrl)
	_, err = a.http.PostForm(AnkiLoginUrl,
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

	resp, _ = a.http.Get(AnkiEditorUrl)
	body, _ := ioutil.ReadAll(resp.Body)
	midMatch := reMid.FindSubmatch(body)

	if len(midMatch) < 2 {
		return errors.New(
			"failed to get mid authentification value from anki web")
	}

	a.mid = string(midMatch[1])

	return nil
}

func (a *AnkiAccount) Search(searchText string) (bool, error) {
	resp, err := a.http.PostForm(AnkiSearchUrl,
		url.Values{
			"keyword":   {searchText},
			"submitted": {"1"},
		})

	if err != nil {
		return false, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	items := reItem.FindAllSubmatch(body, -1)
	for _, m := range items {
		if strings.TrimSpace(string(m[1])) == searchText {
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

	resp, err := a.http.PostForm(AnkiAddUrl, urlValues)

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
