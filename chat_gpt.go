package main

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/reconquest/karma-go"
)

func requestCardFromChatGpt(
	apiKey string,
	term string,
) (string, string, error) {
	client := resty.New()

	response, err := client.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model": "gpt-4",
			"messages": []interface{}{
				map[string]interface{}{
					"role":    "system",
					"content": "You are a polish teacher that helps to build a card deck for a spactial learning.",
				},
				map[string]interface{}{
					"role":    "user",
					"content": getChatGptRequest(term),
				},
			},
		}).
		Post("https://api.openai.com/v1/chat/completions")

	if err != nil {
		return "", "", karma.Format(err, "failed to get response from chat gpt")
	}

	body := response.Body()

	if response.StatusCode() != 200 {
		return "", "", karma.
			Describe("response", string(body)).
			Reason(errors.New("chat gpt responded with error"))
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", "", karma.Format(
			err,
			"failed to decode json response from chat gpt",
		)
	}

	// Extract the content from the JSON response
	content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	parts := strings.Split(content, "Translation:")
	if len(parts) != 2 {
		return "", "", karma.
			Describe("content", content).
			Reason("chat gpt reponse does not contain translation section")

	}

	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
}

func getChatGptRequest(term string) string {
	return `Given a word "` + term + `", return a response that contains 4 sections:

1. The basic form for this word.
2. Three examples of its usage in different forms and different contexts.
3. Word "Translation" (this one is important to make the format machine-readable).
4. The english translation for this word (but not for the examples). The english
   translation may contain as many definitions as needed.

Here is an example how it should be done for another word "pracujemy":

---

Pracować

1. Jestem dobrą osobą i ciężko pracuję.
2. To na pewno wysoko postawiony pracownik.
3. Jeśli pracujesz sam, ubieraj maskę.

Translation:

1. to work, to labor (to exert oneself very much)
2. (intransitive or reflexive with się) to work (to perform a specified task)

---

Do the same for the word "` + term + `".`
}
