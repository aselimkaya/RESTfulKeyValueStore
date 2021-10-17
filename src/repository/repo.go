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

//Entry holds the key value pair
type Entry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//GetStore is a getter function for key value map
func GetStore() map[string]string {
	return keyValStore
}

//Init is initial funciton that sets up the JSON file according to given path
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

//ConverFromJSON takes HTTP request body as parameter and converts it to Entry struct
func (e *Entry) ConvertFromJSON(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	return decoder.Decode(e)
}

//AddEntry function takes key and value as parameters and adds into key value map. This function also checks if given key exists and returns the result.
func AddEntry(key, value string, l *log.Logger) bool {
	isExists := false
	if _, ok := keyValStore[key]; ok {
		l.Println("Key already exists in the store, value will be updated")
		isExists = true
	}
	keyValStore[key] = value
	return isExists
}

//ErrKey used for key not found errors
var ErrKey = fmt.Errorf("key not found")

func GetEntry(key string) (Entry, error) {
	if val, ok := keyValStore[key]; ok {
		return Entry{Key: key, Value: val}, nil
	}
	return Entry{}, ErrKey
}

//Sync writes map content to JSON file
func Sync(filePath string, l *log.Logger) error {
	return utils.SyncFile(filePath, l, keyValStore)
}

//Flush clears not only map but also JSON file content
func Flush(filePath string) error {
	keyValStore = make(map[string]string)
	return utils.FlushFile(filePath)
}
