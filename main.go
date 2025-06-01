package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	mdParser "github.com/gomarkdown/markdown/parser"
	"github.com/olexsmir/anpi/anki"
	"github.com/olexsmir/anpi/parser"
)

//nolint:gochecknoglobals // comment to make linter happy
var cli struct {
	File string `help:"path to import file" name:"file" type:"path"`
}

func main() {
	_ = kong.Parse(&cli)
	if err := run(cli.File); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(fpath string) error {
	ac := anki.NewAnkiClient()

	f, err := os.ReadFile(filepath.Clean(fpath))
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
			nid, err := ac.AddNote(anki.Note{
				DeckName:  deck.Deck,
				ModelName: deck.Type,
				Fields:    fields,
				Tags:      tags,
			})
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
	str = strings.TrimSuffix(str, "\n")

	return str
}
