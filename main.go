package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Sense struct {
	Definition string
	Examples   []string
}

type POSBlock struct {
	POS    string
	IPA    string
	Senses []Sense
}

type Entry struct {
	Word         string
	RedirectWord string
	POSBlocks    []POSBlock
}

func isPiped() bool {
	fi, _ := os.Stdin.Stat()
	return (fi.Mode() & os.ModeCharDevice) == 0
}

func readStdinWord() (string, error) {
	s := bufio.NewScanner(os.Stdin)
	if s.Scan() {
		return strings.TrimSpace(s.Text()), nil
	}
	return "", s.Err()
}

func main() {
	all := flag.Bool("all", false, "print all senses to stderr")
	aFlag := flag.Bool("a", false, "shorthand for --all")
	pos := flag.String("pos", "", "filter by part of speech")
	ipaFlag := flag.Bool("ipa", false, "print only IPA pronunciation")
	clearCacheFlag := flag.Bool("clear-cache", false, "clear cache for word")

	flag.Parse()

	args := flag.Args()

	var word string
	switch {
	case len(args) > 0:
		word = args[0]
	case isPiped():
		var err error
		word, err = readStdinWord()
		if err != nil {
			fmt.Fprintf(os.Stderr, "anpi: %v\n", err)
			os.Exit(1)
		}
		if word == "" {
			fmt.Fprintln(os.Stderr, "usage: anpi [flags] <word>")
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "usage: anpi [flags] <word>")
		os.Exit(1)
	}

	if len(args) > 1 {
		for _, arg := range args[1:] {
			if arg == "--all" || arg == "-a" {
				*all = true
			}
			if arg == "--clear-cache" || arg == "--ipa" {
				*clearCacheFlag = true
			}
		}
	}

	if *all || *aFlag {
		*all = true
	}

	cache, err := NewCache()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *clearCacheFlag {
		_ = cache.Clear(word)
	}

	entry, err := cache.Read(word)
	if errors.Is(err, ErrNotFound) {
		_, _ = fmt.Fprintf(os.Stderr, "anpi: %v\n", err)
		os.Exit(1)
	}
	if err != nil {
		dict := NewCambridgeDictionary()
		entry, err = dict.Lookup(word)
		if err != nil {
			if werr := cache.Write(word, nil, false); werr != nil {
				fmt.Fprintf(os.Stderr, "cache write error: %v\n", werr)
			}
			_, _ = fmt.Fprintf(os.Stderr, "anpi: %v\n", err)
			os.Exit(1)
		}
		if err := cache.Write(word, entry, true); err != nil {
			fmt.Fprintf(os.Stderr, "cache write error: %v\n", err)
		}
	}

	if *pos != "" {
		var filtered []POSBlock
		for _, b := range entry.POSBlocks {
			if strings.EqualFold(b.POS, *pos) {
				filtered = append(filtered, b)
			}
		}
		if len(filtered) == 0 {
			_, _ = fmt.Fprintf(os.Stderr, "warning: no %q sense found, showing all\n", *pos)
		} else {
			entry.POSBlocks = filtered
		}
	}

	if *ipaFlag {
		for _, b := range entry.POSBlocks {
			if b.IPA != "" {
				fmt.Printf("/%s/\n", b.IPA)
			}
		}
		return
	}

	if *all {
		writeAll(entry, os.Stderr)
		return
	}

	fmt.Print(toTSV(entry))
}
