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
	Meanings   []string
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

	lang := flags.String("lang", "en-ru",
		"translation direction")
	creds := flags.String("creds", os.Getenv("HOME")+"/.runki/creds",
		"path to creds file")
	user := flags.String("user", "",
		"ankiweb username")
	pass := flags.String("pass", "",
		"ankiweb password")
	deck := flags.String("deck", "english",
		"deck to add")
	dry := flags.Bool("dry", false,
		"dry run (do not alter anki db)")
	cut := flags.Int("cut", 0,
		"stop processing after N non-unique words found")
	silent := flags.Bool("s", false,
		"silent, do not print translation before uniq check")

	conf := loadConfig(os.Getenv("HOME") + "/.runki/runkirc")

	flags.Parse(append(conf, os.Args[1:]...))

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
			if !*silent {
				fmt.Println("<" + unknown + ": no translation found>")
			}

			continue
		}

		translation := "[" + lookup.Transcript + "] " +
			strings.Join(lookup.Meanings, ", ")

		if !*silent {
			fmt.Println(translation)
		}

		if *dry {
			continue
		}

		found, err := anki.Search(unknown)
		if err != nil {
			log.Fatalf(err.Error())
		}

		if !found {
			if *silent {
				fmt.Println(translation)
			}

			foundStreak = 0
			err = anki.Add(*deck, unknown, translation)
			if err != nil {
				log.Fatalln(err.Error())
			}
		} else {
			foundStreak += 1
		}

		if foundStreak >= *cut && *cut > 0 {
			log.Fatalf("stopping after %d consequent non-unique words", *cut)
			break
		}
	}
}