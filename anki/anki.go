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
