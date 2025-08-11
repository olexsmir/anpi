package main

import (
	"os"
	"strings"
	"testing"

	"olexsmir.xyz/x/is"
)

func TestCambridgeDictionary_Scrape_skip(t *testing.T) {
	f, err := os.Open("testdata/cambridge_skip.html")
	is.Err(t, err, nil)
	defer func() { _ = f.Close() }()

	dict := NewCambridgeDictionary()
	entry, err := dict.Scrape("skip", f)
	is.Err(t, err, nil)
	is.Equal(t, "skip", entry.Word)
	is.Equal(t, true, len(entry.POSBlocks) > 0)

	block := entry.POSBlocks[0]
	is.Equal(t, true, block.POS != "")
	is.Equal(t, true, strings.HasPrefix(block.IPA, "sk") || strings.HasPrefix(block.IPA, "/sk"))

	is.Equal(t, true, len(block.Senses) > 0)
	is.Equal(t, true, block.Senses[0].Definition != "")

	is.Equal(t, true, len(block.Senses[0].Examples) > 0)
	is.Equal(t, true, strings.Contains(block.Senses[0].Examples[0], "skip"))
}

func TestCambridgeDictionary_Scrape_hello(t *testing.T) {
	f, err := os.Open("testdata/cambridge_hello.html")
	is.Err(t, err, nil)
	defer func() { _ = f.Close() }()

	dict := NewCambridgeDictionary()
	entry, err := dict.Scrape("hello", f)
	is.Err(t, err, nil)
	is.Equal(t, "hello", entry.Word)
	is.Equal(t, true, len(entry.POSBlocks) > 0)

	block := entry.POSBlocks[0]
	is.Equal(t, true, block.POS != "")
	is.Equal(t, true, block.IPA != "")
	is.Equal(t, true, len(block.Senses) > 0)
	is.Equal(t, true, block.Senses[0].Definition != "")
}

func TestCambridgeDictionary_Scrape_notFound(t *testing.T) {
	f, err := os.Open("testdata/cambridge_notfound.html")
	is.Err(t, err, nil)
	defer func() { _ = f.Close() }()

	dict := NewCambridgeDictionary()
	_, err = dict.Scrape("asdfghjklzxcvbnm", f)
	is.Err(t, err, "word not found")
}

func TestCambridgeDictionary_Scrape_skipHasExamples(t *testing.T) {
	f, err := os.Open("testdata/cambridge_skip.html")
	is.Err(t, err, nil)
	defer func() { _ = f.Close() }()

	dict := NewCambridgeDictionary()
	entry, err := dict.Scrape("skip", f)
	is.Err(t, err, nil)

	is.Equal(t, true, len(entry.POSBlocks[0].Senses[0].Examples) > 0)
	is.Equal(t, true, len(entry.POSBlocks[0].Senses[0].Examples) <= 3)
}

func TestCambridgeDictionary_Scrape_helloHasExamples(t *testing.T) {
	f, err := os.Open("testdata/cambridge_hello.html")
	is.Err(t, err, nil)
	defer func() { _ = f.Close() }()

	dict := NewCambridgeDictionary()
	entry, err := dict.Scrape("hello", f)
	is.Err(t, err, nil)

	exs := entry.POSBlocks[0].Senses[0].Examples
	is.Equal(t, true, len(exs) > 0)
	is.Equal(t, true, len(exs) <= 3)
}

func TestCambridgeDictionary_Scrape_helloPOS(t *testing.T) {
	f, err := os.Open("testdata/cambridge_hello.html")
	is.Err(t, err, nil)
	defer func() { _ = f.Close() }()

	dict := NewCambridgeDictionary()
	entry, err := dict.Scrape("hello", f)
	is.Err(t, err, nil)

	is.Equal(t, "exclamation", entry.POSBlocks[0].POS)
}

func TestCambridgeDictionary_Scrape_multiplePOSBlocks(t *testing.T) {
	f, err := os.Open("testdata/cambridge_skip.html")
	is.Err(t, err, nil)
	defer func() { _ = f.Close() }()

	dict := NewCambridgeDictionary()
	entry, err := dict.Scrape("skip", f)
	is.Err(t, err, nil)

	is.Equal(t, true, len(entry.POSBlocks) > 1)
}

func TestCambridgeDictionary_Scrape_multipleSenses(t *testing.T) {
	f, err := os.Open("testdata/cambridge_skip.html")
	is.Err(t, err, nil)
	defer func() { _ = f.Close() }()

	dict := NewCambridgeDictionary()
	entry, err := dict.Scrape("skip", f)
	is.Err(t, err, nil)

	is.Equal(t, true, len(entry.POSBlocks[0].Senses) > 1)
}

func TestCambridgeDictionary_Scrape_usIPA(t *testing.T) {
	f, err := os.Open("testdata/cambridge_skip.html")
	is.Err(t, err, nil)
	defer func() { _ = f.Close() }()

	dict := NewCambridgeDictionary()
	entry, err := dict.Scrape("skip", f)
	is.Err(t, err, nil)

	is.Equal(t, true, entry.POSBlocks[0].IPA != "")
}
