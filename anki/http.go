package anki

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ankiURL            = "http://localhost:8765"
	ankiConnectVersion = 6
)

type AnkiResponse[T any] struct {
	Result T      `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func (r AnkiResponse[T]) CheckErrors() error {
	if r.Error != "" {
		return fmt.Errorf("%s", r.Error)
	}
	return nil
}

type AnkiRequest[T any] struct {
	Action  string `json:"action"`
	Version int    `json:"version"`
	Params  T      `json:"params"`
}

type paramsDefault struct{}

func request[R any, P any](action string, params P) (AnkiResponse[R], error) {
	bodyReq, err := json.Marshal(AnkiRequest[P]{
		Action:  action,
		Version: ankiConnectVersion,
		Params:  params,
	})
	if err != nil {
		return AnkiResponse[R]{}, err
	}

	req, err := http.NewRequest(http.MethodGet, ankiURL, bytes.NewBuffer(bodyReq))
	if err != nil {
		return AnkiResponse[R]{}, err
	}

	req.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return AnkiResponse[R]{}, err
	}

	defer resp.Body.Close() //nolint

	var ankiResp AnkiResponse[R]
	if err := json.NewDecoder(resp.Body).Decode(&ankiResp); err != nil {
		return AnkiResponse[R]{}, err
	}

	return ankiResp, ankiResp.CheckErrors()
}
