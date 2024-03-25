package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type LookupResult struct {
	Transcript string    `json:"transcript"`
	Meanings   []Meaning `json:"meanings"`
}

type Meaning struct {
	Translation string   `json:"translation"`
	References  []string `json:"references"`
}

func (meaning *Meaning) String() string {
	return fmt.Sprintf(
		`%s (%s)`,
		meaning.Translation,
		strings.Join(meaning.References, `, `),
	)
}

func loadConfig(path string) []string {
	args := make([]string, 0)

	conf, err := ioutil.ReadFile(path)
	if err != nil {
		return args
	}

	confLines := strings.Split(string(conf), "\n")
	for _, line := range confLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		args = append(args, line)
	}

	return args
}

func main() {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flags.Usage = displayHelp

	lang := flags.String("lang", "en-ru", "translation direction")
	creds := flags.String("creds", os.Getenv("HOME")+"/.config/runki/creds",
		"path to creds file")
	chatgptApiKey := flags.String("chatgpt-api-key", "", "chat gpt api key")
	// user := flags.String("user", "", "ankiweb username")
	// pass := flags.String("pass", "", "ankiweb password")
	deckId := flags.Int64("deck-id", 0, "deck to add")
	notetypeId := flags.Int64("notetype-id", 0, "deck to add")
	dry := flags.Bool("dry", false, "dry run (do not alter anki db)")
	cut := flags.Int("cut", 0, "stop processing after N non-unique words found")
	silent := flags.Bool("silent", false, "silent, do not print translation "+
		"before uniq check")
	useJSON := flags.Bool("json", false, "output in json")

	conf := loadConfig(os.Getenv("HOME") + "/.config/runki/runkirc")

	flags.Parse(append(conf, os.Args[1:]...))

	addCard(
		*lang,
		*creds,
		*chatgptApiKey,
		*deckId,
		*notetypeId,
		*dry,
		*cut,
		*silent,
		*useJSON,
	)
}

func displayHelp() {
	fmt.Println(`
NAME
	anki - ankiweb and yandex-dictionary client. Provides cli interface for
	adding word and translation to http://ankiweb.net.

SYNOPSIS
	runki [--lang LANG] [--creds CREDS] [--user USER] [--pass PASS]
		[--deck DECK] [--dry] [--cut CUT] [--silent] [--help]

DESCRIPTION
	--lang
		Translation direction; [from]-[to], for example en-ru; see
		http://api.yandex.com/dictionary/doc/dg/reference/lookup.xml
		for detailed description. Default: en-ru.

	--creds
		The cookie stolen from the ankiweb session.

	--chatgpt-api-key
		ChatGPT API key.

	--deck-id
		Id of the deck to add the card to. To find out this value, add the card via
		the web ui, check the stack trace of the network request, set breakpoint to
		stop at the request time, and then, debug the value inside the code. The
		process is tricky and the improvements for it are welcome.

	--notetype-id
		Id of the notetype to add the card to. Can be found in the same way as the
		deck id.

	--dry
		Do not add card, just show translation.

	--cut
		Specifies how many duplications can be found before stopping input
		processing. If zero - dont process without stops. Default: 0

	--silent
		Do not print messages.

	--help
		Display this help.

	All options can be red from ~/.config/runki/runkirc.

EXAMPLES
	# add test card with default settings
	echo test | runki

	# add test card with user and password
	echo test | runki --user user@example.com --pass PASSWORD

	# ~/.config/runki/runkirc
	--creds
	has_auth=1; ankiweb=eyJvcCI6ImNrIiwiaWF0IjoxNzA4ODA2XXX

SEE ALSO
	http://ankiweb.net
	http://api.yandex.com/dictionary/
	https://github.com/seletskiy/runki

AUTHORS
	Stanislav Seletskiy
	Leonid Shagabutdinov

VERSION
	2.0
`)
}

func addCard(
	lang string,
	creds string,
	chatgptApiKey string,
	deckId int64,
	notetypeId int64,
	dry bool,
	cut int,
	silent bool,
	useJSON bool,
) {
	anki := NewAnkiAccount()

	// shouldAuth, err := anki.Load(creds)
	// if err != nil {
	// 	log.Fatalf("can't read from creds file:", err)
	// }

	// if !dry && shouldAuth {
	// 	err := anki.WebLogin(user, pass)
	// 	if err != nil {
	// 		log.Fatalf("can't login to ankiweb", err.Error())
	// 	}
	// }

	// err = anki.Save(creds)
	// if err != nil {
	// 	log.Fatalf("can't save creds file:", err.Error())
	// 	return
	// }

	stdin := bufio.NewReader(os.Stdin)
	for {
		line, err := stdin.ReadString('\n')
		if err != nil {
			break
		}

		unknown, translation, err := requestCardFromChatGpt(
			chatgptApiKey,
			strings.TrimSpace(line),
		)

		if err != nil {
			log.Fatalf("failed to get translation from chat gpt", err.Error())
		}

		if dry {
			continue
		}

		if silent && !useJSON {
			fmt.Println(translation)
		}

		err = anki.Add(creds, deckId, notetypeId, unknown, translation)
		if err != nil {
			log.Fatal(err)
		}
	}
}
