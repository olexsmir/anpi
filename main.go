package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	mdParser "github.com/gomarkdown/markdown/parser"
	"github.com/olexsmir/anpi/anki"
	"github.com/olexsmir/anpi/parser"
)

const importFile = "import.yml"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	anki := anki.NewAnkiClient()

	f, err := os.ReadFile(importFile)
	if err != nil {
		return err
	}

	data, err := parser.Parse(f)
	if err != nil {
		return err
	}

	for _, deck := range data {
		for _, note := range deck.Notes {
			fields := make(map[string]string)
			for k, v := range note.Fields {
				fields[deck.FieldLookUp(k)] = fromMarkdown(v)
			}

			slog.Info("gotten fields", "fields", fields)

			tags := mergeTags(deck.Tags, note.Tags)
			nid, err := anki.AddNote(deck.Deck, deck.Type, fields, tags)
			if err != nil {
				return err
			}

			slog.Info("note added", "id", nid, "fields", fields)
		}
	}

	return nil
}

func mergeTags(global, local []string) []string {
	unique := make(map[string]struct{})

	for _, tag := range global {
		unique[tag] = struct{}{}
	}
	for _, tag := range local {
		unique[tag] = struct{}{}
	}

	result := make([]string, 0, len(unique))
	for tag := range unique {
		result = append(result, tag)
	}

	return result
}

func fromMarkdown(inp string) string {
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}

	p := mdParser.New()
	doc := p.Parse([]byte(inp))

	str := string(markdown.Render(doc, html.NewRenderer(opts)))
	str = strings.ReplaceAll(str, "<p>", "")
	str = strings.ReplaceAll(str, "</p>", "")

	return str
}
