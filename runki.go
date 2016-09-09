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
	Transcript string
	Meanings   []Meaning
}

type Meaning struct {
	Translation string
	References  []string
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
	user := flags.String("user", "", "ankiweb username")
	pass := flags.String("pass", "", "ankiweb password")
	deck := flags.String("deck", "Default", "deck to add")
	dry := flags.Bool("dry", false, "dry run (do not alter anki db)")
	cut := flags.Int("cut", 0, "stop processing after N non-unique words found")
	silent := flags.Bool("silent", false, "silent, do not print translation "+
		"before uniq check")

	conf := loadConfig(os.Getenv("HOME") + "/.config/runki/runkirc")

	flags.Parse(append(conf, os.Args[1:]...))

	addCard(*lang, *creds, *user, *pass, *deck, *dry, *cut, *silent)
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
		Creds file to cache cookies to. If you changed password or in case of
		authentification failure, delete this file. Default: ~/.config/runki/creds.

	--user
		Your username to http://ankiweb.net.

	--pass
		Your password to http://ankiweb.net.

	--deck
		To add card to. Default: Default.

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
	--user
		user@example.com

	--pass
		PASSWORD

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

func addCard(lang string, creds string, user string, pass string,
	deck string, dry bool, cut int, silent bool) {

	ya := NewYandexProvider(lang, "", UnlimitedSynonyms)
	anki := NewAnkiAccount()

	shouldAuth, err := anki.Load(creds)
	if err != nil {
		log.Fatalf("can't read from creds file:", err)
	}

	if !dry && shouldAuth {
		err := anki.WebLogin(user, pass)
		if err != nil {
			log.Fatalf("can't login to ankiweb", err.Error())
		}
	}

	err = anki.Save(creds)
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
			if !silent {
				fmt.Fprintf(os.Stderr,
					"<"+unknown+": no translation found>")
			}

			continue
		}

		meanings := []string{}
		for _, meaning := range lookup.Meanings {
			meanings = append(meanings, meaning.String())
		}

		translation := ""
		if lookup.Transcript != "" {
			translation = "[" + lookup.Transcript + "] "
		}

		translation = translation + strings.Join(meanings, ", ")

		if !silent {
			fmt.Println(translation)
		}

		if dry {
			continue
		}

		found, err := anki.Search(unknown)
		if err != nil {
			log.Fatalf(err.Error())
		}

		if !found {
			if silent {
				fmt.Println(translation)
			}

			foundStreak = 0
			err = anki.Add(deck, unknown, translation)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			foundStreak += 1
		}

		if foundStreak >= cut && cut > 0 {
			log.Fatalf("stopping after %d consequent non-unique words", cut)
		}
	}
}
