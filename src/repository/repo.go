package repository

import (
	"encoding/json"
	"io"
)

var KeyValStore map[string]string

type Entry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (e *Entry) ConvertFromJSON(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(e)
}

func AddPair(key, value string) bool {
	KeyValStore[key] = value
	return true
}
