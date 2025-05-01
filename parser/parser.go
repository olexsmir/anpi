package parser

import (
	"errors"

	"gopkg.in/yaml.v3"
)

type DeckImport struct {
	Deck   string            `yaml:"deck"`
	Type   string            `yaml:"type"`
	Tags   []string          `yaml:"tags"`
	Fields map[string]string `yaml:"fields"`
	Notes  []Note            `yaml:"notes"`
}

func (d *DeckImport) FieldLookUp(field string) string {
	// TODO: use ok notation here

	val := d.Fields[field]
	return val
}

func (d *DeckImport) Validate() error {
	return nil
}

type Note struct {
	Fields map[string]string
	Tags   []string
}

func (n *Note) UnmarshalYAML(value *yaml.Node) error {
	var raw map[string]any
	if err := value.Decode(&raw); err != nil {
		return err
	}

	n.Fields = make(map[string]string)
	for k, v := range raw {
		if k == "tags" {
			if tags, ok := v.([]string); ok {
				n.Tags = append(n.Tags, tags...)
			}
		}

		if str, ok := v.(string); ok {
			n.Fields[k] = str
		}
	}

	return nil
}

var ErrInvalidFileFormat = errors.New("invalid file format")

func Parse(inp []byte) ([]DeckImport, error) {
	var listRes []DeckImport
	if err := yaml.Unmarshal(inp, &listRes); err == nil {
		return listRes, nil
	}

	var single DeckImport
	if err := yaml.Unmarshal(inp, &single); err == nil {
		return []DeckImport{single}, nil
	}

	return nil, ErrInvalidFileFormat
}
