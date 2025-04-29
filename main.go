package main

import (
	"log/slog"

	"github.com/olexsmir/anpi/anki"
)

func main() {
	anki := anki.NewAnkiClient()

	nid, err := anki.AddNote("testing", "Basic", map[string]string{
		"Front": "the thing",
		"Back":  "asdfasfasdf",
	})
	slog.Info("note added", "nid", nid, "err", err)

	// decks, err := anki.GetDeckNames()
	// slog.Info("test deck names", "fields", decks, "err", err)
	//
	// fields, err := anki.GetModelFieldNames("Basic")
	// slog.Info("test existent type", "fields", fields, "err", err)
	//
	// fieldsErr, err := anki.GetModelFieldNames("chicked jokey")
	// slog.Info("test non-existent type", "fields", fieldsErr, "err", err)
	//
	// f, _ := os.ReadFile("test.yml")
	// data, err := parser.Parse(f)
	// slog.Info("parsing data", "err", err)
	//
	// lookup := "front"
	// ankiField := data[1].FieldLookUp(lookup)
	// slog.Info("looking up anki note field", "field", ankiField, "lookup", lookup)

	// for _, deck := range data {
	// 	fmt.Printf("\nDeck: %s\n", deck.Deck)
	// 	fmt.Printf("Type: %s\n", deck.Type)
	// 	fmt.Printf("Global Tags: %v\n", deck.Tags)
	// 	fmt.Println("Fields mapping:")
	// 	for internalName, anki := range deck.Fields {
	// 		fmt.Printf("  %s -> %s\n", internalName, anki)
	// 	}
	//
	// 	fmt.Println("Notes:")
	// 	for _, note := range deck.Notes {
	// 		fmt.Println("  Note fields:")
	// 		for internalName, value := range note.Fields {
	// 			fmt.Printf("    %s: %s\n", internalName, value)
	// 		}
	// 		if len(note.Tags) > 0 {
	// 			fmt.Printf("    Local Tags: %v\n", note.Tags)
	// 		}
	// 	}
	// }
}
