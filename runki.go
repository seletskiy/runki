package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type LookupResult struct {
	Transcript string
	Meanings   []string
}

var lang = flag.String("lang", "en-ru", "translation direction")
var creds = flag.String("creds", os.Getenv("HOME")+"/.runki/creds",
	"path to *creds file")
var user = flag.String("user", "", "ankiweb username")
var pass = flag.String("pass", "", "ankiweb password")
var deck = flag.String("deck", "english", "deck to add")

func main() {
	flag.Parse()

	ya := NewYandexProvider(*lang, "")
	anki := NewAnkiAccount()

	err := anki.Load(*creds)
	if err != nil {
		log.Println("can't read from creds file:", err)
		err := anki.WebLogin(*user, *pass)
		if err != nil {
			log.Fatalf("can't login to ankiweb", err.Error())
		}
	}

	err = anki.Save(*creds)
	if err != nil {
		log.Fatalf("can't save creds file:", err.Error())
	}

	stdin := bufio.NewReader(os.Stdin)
	for {
		line, err := stdin.ReadString('\n')
		if err != nil {
			break
		}

		unknown := strings.TrimSpace(line)

		lookup, err := ya.Lookup(unknown)
		if err != nil {
			log.Fatalf(err.Error())
		}

		translation := "[" + lookup.Transcript + "] " +
			strings.Join(lookup.Meanings, ", ")

		fmt.Println(translation)

		found, err := anki.Search(unknown)
		if err != nil {
			log.Fatalf(err.Error())
		}

		if !found {
			anki.Add(*deck, unknown, translation)
		}
	}

	return
}
