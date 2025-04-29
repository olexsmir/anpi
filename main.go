package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/olexsmir/anpi/anki"
	"github.com/olexsmir/anpi/parser"
)

func main() {
	anki := anki.NewAnkiClient()

	decks, err := anki.GetDeckNames()
	slog.Info("test deck names", "fields", decks, "err", err)

	fields, err := anki.GetModelFieldNames("Basic")
	slog.Info("test existent type", "fields", fields, "err", err)

	fieldsErr, err := anki.GetModelFieldNames("chicked jokey")
	slog.Info("test non-existent type", "fields", fieldsErr, "err", err)

	f, _ := os.ReadFile("test.yml")
	data, err := parser.Parse(f)
	slog.Info("parsing data", "err", err)

	for _, deck := range data {
		fmt.Printf("\nDeck: %s\n", deck.Deck)
		fmt.Printf("Type: %s\n", deck.Type)
		fmt.Printf("Global Tags: %v\n", deck.Tags)
		fmt.Println("Fields mapping:")
		for internalName, anki := range deck.Fields {
			fmt.Printf("  %s -> %s\n", internalName, anki)
		}

		fmt.Println("Notes:")
		for _, note := range deck.Notes {
			fmt.Println("  Note fields:")
			for internalName, value := range note.Fields {
				fmt.Printf("    %s: %s\n", internalName, value)
			}
			if len(note.Tags) > 0 {
				fmt.Printf("    Local Tags: %v\n", note.Tags)
			}
		}
	}
}
