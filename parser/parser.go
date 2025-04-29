package parser

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type DeckImport struct {
	Deck   string            `yaml:"deck"`
	Type   string            `yaml:"type"`
	Tags   []string          `yaml:"tags"`
	Fields map[string]string `yaml:"fields"`
	Notes  []Note            `yaml:"notes"`

	// fieldLookUp internal reserve mapping of [Fields]
	fieldLookUp map[string]string
}

func (d *DeckImport) FieldLookUp(field string) string {
	// TODO: use ok notation here

	val := d.fieldLookUp[field]
	return val
}

func (d *DeckImport) Validate() error {
	return nil
}

func (d *DeckImport) buildLookUpTable() {
	d.fieldLookUp = make(map[string]string)
	for k, v := range d.Fields {
		d.fieldLookUp[k] = v
	}
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

func Parse(inp []byte) ([]DeckImport, error) {
	var listRes []DeckImport
	if err := yaml.Unmarshal(inp, &listRes); err == nil {
		for i := range listRes {
			listRes[i].buildLookUpTable()
		}
		return listRes, nil
	}

	var single DeckImport
	if err := yaml.Unmarshal(inp, &single); err == nil {
		single.buildLookUpTable()
		return []DeckImport{single}, nil
	}

	return nil, fmt.Errorf("invalid file format")
}
