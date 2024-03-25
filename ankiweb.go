package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/reconquest/karma-go"
	"github.com/seletskiy/runki/messages"
	"google.golang.org/protobuf/proto"
)

var _ = fmt.Print

type AnkiAccount struct {
	http *http.Client
	mid  string
}

const (
	AnkiBaseUrl        = "https://ankiweb.net"
	AnkiLoginUrl       = AnkiBaseUrl + "/svc/account/login"
	AnkiCheckCookieUrl = AnkiBaseUrl + "/account/checkCookie"
	AnkiSearchUrl      = AnkiBaseUrl + "/search/"
	AnkiAddUrl         = AnkiBaseUrl + "/svc/editor/add-or-update"
	// AnkiEditorUrl      = AnkiBaseUrl + "/edit/"
)

var (
	reItem = regexp.MustCompile(`(?s:mitem3.*?<td>([^/]+))`)
)

func NewAnkiAccount() *AnkiAccount {
	return &AnkiAccount{}
}

func (a *AnkiAccount) Load(filename string) (shouldAuth bool, err error) {
	storedData := struct {
		Cookies []*http.Cookie
		Mid     string
	}{
		make([]*http.Cookie, 0),
		"",
	}

	data, err := ioutil.ReadFile(filename)
	switch {
	case err == nil:
		err = json.Unmarshal(data, &storedData)
		if err != nil {
			return false, err
		}

	case os.IsNotExist(err):
		shouldAuth = true

	default:
		return false, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return false, err
	}

	ankiurl, _ := url.Parse(AnkiBaseUrl)
	jar.SetCookies(ankiurl, storedData.Cookies)

	a.http = &http.Client{Jar: jar}
	a.mid = storedData.Mid

	return shouldAuth, nil
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

func getResponseBodyOrError(response *http.Response) string {
	respBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err.Error()
	} else {
		return string(respBytes)
	}
}

func (a *AnkiAccount) WebLogin(login, password string) error {
	jar, _ := cookiejar.New(nil)
	a.http = &http.Client{Jar: jar}

	getResponse, err := a.http.Get(AnkiBaseUrl)
	if err != nil {
		return errors.New("failed to get reponse from anki web; " + err.Error() +
			getResponseBodyOrError(getResponse))
	}

	postResponse, err := a.http.Post(
		"https://ankiweb.net/svc/account/login",
		"application/octet-stream",
		strings.NewReader("\n\x18"+login+"\x12\n"+password),
	)
	if err != nil {
		return karma.Format(err, "failed to send authorization request")
	}

	if postResponse.StatusCode != 200 {
		return karma.
			Describe("status", postResponse.StatusCode).
			Describe("response", getResponseBodyOrError(postResponse)).
			Reason("failed to login to anki web")
	}

	checkResponse, err := a.http.Post(
		"https://ankiweb.net/svc/account/get-account-status",
		"application/octet-stream",
		strings.NewReader(""),
	)
	if err != nil {
		return karma.Format(err, "failed to validate authorization")
	}

	if checkResponse.StatusCode != 200 {
		return karma.
			Describe("status", checkResponse.StatusCode).
			Describe("response", getResponseBodyOrError(checkResponse)).
			Reason("failed to login to anki web")
	}

	ankiurl, _ := url.Parse("https://ankiweb.net")
	if len(jar.Cookies(ankiurl)) == 0 {
		responseError := ""
		respBytes, err := io.ReadAll(checkResponse.Body)
		if err != nil {
			responseError = err.Error()
		} else {
			responseError = string(respBytes)
		}

		return errors.New("failed to login to anki web (response: " + responseError + ")")
	}

	return nil
}

func (a *AnkiAccount) Search(searchText string) (bool, error) {
	resp, err := a.http.Post(
		AnkiSearchUrl,
		"application/octet-stream",
		strings.NewReader("\x12"+searchText),
	)

	if err != nil {
		return false, err
	}

	body, _ := io.ReadAll(resp.Body)

	items := reItem.FindAllSubmatch(body, -1)
	for _, m := range items {
		if strings.TrimSpace(string(m[1])) == searchText {
			return true, nil
		}
	}

	return false, nil
}

// func printStringWithEscapeChars(bytes []byte) {
// 	for _, char := range bytes {
// 		// Check if the character is printable
// 		if char == '\n' {
// 			fmt.Printf("\\n")
// 		} else if char >= 32 && char <= 126 {
// 			fmt.Printf("%c", char) // Print the printable character as-is
// 		} else {
// 			// Print the character using Unicode escape sequence \xXX
// 			fmt.Printf("\\x%02X", char)
// 		}
// 	}
// 	fmt.Println()
// }

func (a *AnkiAccount) Add(
	cookie string,
	deckId, notetypeId int64,
	text, translation string,
) error {
	message := &messages.Message{
		// Fields: []string{"test", "test"},
		Fields: []string{
			"<div style=\"text-align: left;\">" + strings.ReplaceAll(text, "\n", "<br>") + "</div>",
			"<div style=\"text-align: left;\">" + strings.ReplaceAll(translation, "\n", "<br>") + "</div>",
		},
		Tags: "",
		Mode: &messages.Message_Add{
			Add: &messages.AddOrUpdateRequest_AddMode{
				DeckId:     deckId,
				NotetypeId: notetypeId,
			},
		},
	}

	data, err := proto.Marshal(message)
	if err != nil {
		return karma.Format(err, "failed to marshal data")
	}

	// printStringWithEscapeChars(data)

	req, err := http.NewRequest(
		"POST",
		"https://ankiuser.net/svc/editor/add-or-update",
		strings.NewReader(string(data)),
	)

	if err != nil {
		return karma.Format(err, "failed to create network request")
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Cookie", cookie)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return karma.Format(err, "failed to create anki card")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return karma.
			Describe("status", resp.StatusCode).
			Describe("body", body).
			Reason("ankiweb responded with error")
	}

	return nil
}
