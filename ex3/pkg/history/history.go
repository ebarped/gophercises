package history

import (
	"encoding/json"
	"os"
)

type Options []struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options Options  `json:"options"`
}

type History map[string]Chapter

func New(filename string) (History, error) {
	h, err := parseHistory(filename)
	if err != nil {
		return nil, err
	}
	return h, nil

}

func parseHistory(filename string) (History, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// create a new decoder
	d := json.NewDecoder(f)
	if err != nil {
		return nil, err
	}

	// initialize the storage for the decoded data
	var history History

	// decode the data
	err = d.Decode(&history)
	if err != nil {
		return nil, err
	}
	return history, nil
}
