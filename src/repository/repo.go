package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/aselimkaya/RESTfulKeyValueStore/src/utils"
)

var keyValStore map[string]string

type Entry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func Init(l *log.Logger, path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			l.Fatal(err)
			return
		}

		_, err = f.WriteString(`{}`)
		if err != nil {
			l.Fatal(err)
			return
		}

		keyValStore = make(map[string]string)
	} else {
		jsonFile, err := os.Open(path)
		if err != nil {
			l.Fatal(err)
		}
		defer jsonFile.Close()

		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			l.Fatal(err)
			return
		}

		json.Unmarshal([]byte(byteValue), &keyValStore)
	}
}

func (e *Entry) ConvertFromJSON(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(e)
}

func AddEntry(key, value string, l *log.Logger) bool {
	isExists := false
	if _, ok := keyValStore[key]; ok {
		l.Println("Key already exists in the store, value will be updated")
		isExists = true
	}
	keyValStore[key] = value
	return isExists
}

var ErrKey = fmt.Errorf("key not found")

func GetEntry(key string) (Entry, error) {
	if val, ok := keyValStore[key]; ok {
		return Entry{Key: key, Value: val}, nil
	}
	return Entry{}, ErrKey
}

func Sync(filePath string, l *log.Logger) error {
	return utils.SyncFile(filePath, l, keyValStore)
}

func Flush(filePath string) error {
	keyValStore = make(map[string]string)
	return utils.FlushFile(filePath)
}
