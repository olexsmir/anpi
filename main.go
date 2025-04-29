package main

import (
	"log/slog"

	"github.com/olexsmir/anpi/anki"
)

func main() {
	anki := anki.NewAnkiClient()

	decks, err := anki.GetDeckNames()
	slog.Info("test deck names", "fields", decks, "err", err)

	fields, err := anki.GetModelFieldNames("Basic")
	slog.Info("test existent type", "fields", fields, "err", err)

	fieldsErr, err := anki.GetModelFieldNames("chicked jokey")
	slog.Info("test non-existent type", "fields", fieldsErr, "err", err)
}
