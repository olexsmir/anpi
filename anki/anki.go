package anki

type AnkiClient struct{}

func NewAnkiClient() AnkiClient {
	return AnkiClient{}
}

type ModelFieldNames []string

type ModelFieldNamesParams struct {
	ModelName string `json:"modelName"`
}

func (a *AnkiClient) GetModelFieldNames(modelName string) (ModelFieldNames, error) {
	res, err := request[ModelFieldNames]("modelFieldNames", ModelFieldNamesParams{
		ModelName: modelName,
	})
	if err != nil {
		return ModelFieldNames{}, err
	}

	return res.Result, nil
}

type DeckNames []string

func (a *AnkiClient) GetDeckNames() (DeckNames, error) {
	res, err := request[DeckNames]("deckNames", paramsDefault{})
	if err != nil {
		return DeckNames{}, err
	}

	return res.Result, nil
}

type (
	NoteID int
	Note   struct {
		DeckName  string            `json:"deckName"`
		ModelName string            `json:"modelName"`
		Fields    map[string]string `json:"fields"`
		Tags      []string          `json:"tags"`
	}

	internalNote struct {
		Note Note `json:"note"`
	}
)

func (a *AnkiClient) AddNote(inp Note) (NoteID, error) {
	res, err := request[NoteID]("addNote", internalNote{Note{
		DeckName:  inp.DeckName,
		ModelName: inp.ModelName,
		Fields:    inp.Fields,
		Tags:      inp.Tags,
	}})
	if err != nil {
		return 0, err
	}

	return res.Result, nil
}
