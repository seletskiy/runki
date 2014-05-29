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
	"path to creds file")
var user = flag.String("user", "", "ankiweb username")
var pass = flag.String("pass", "", "ankiweb password")
var deck = flag.String("deck", "english", "deck to add")
var dry = flag.Bool("dry", false, "dry run (do not alter anki db)")
var cut = flag.Int("cut", 0, "stop processing after N non-unique words found")

func main() {
	flag.Parse()

	ya := NewYandexProvider(*lang, "", 2)
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
	foundStreak := 0
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

		if lookup == nil {
			fmt.Println("<" + unknown + ": no translation found>")
			continue
		}

		translation := "[" + lookup.Transcript + "] " +
			strings.Join(lookup.Meanings, ", ")

		fmt.Println(translation)

		if *dry {
			continue
		}

		found, err := anki.Search(unknown)
		if err != nil {
			log.Fatalf(err.Error())
		}

		if !found {
			foundStreak = 0
			anki.Add(*deck, unknown, translation)
		} else {
			foundStreak += 1
		}

		if foundStreak >= *cut && *cut > 0 {
			log.Fatalf("stopping after %d consequent non-unique words", *cut)
			break
		}
	}

	return
}
