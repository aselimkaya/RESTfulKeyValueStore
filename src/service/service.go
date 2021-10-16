package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aselimkaya/RESTfulKeyValueStore/src/repository"
)

type Store struct {
	storeLogger  *log.Logger
	jsonFilePath string
}

func New(l *log.Logger, path string) *Store {
	repository.Init(l, path)
	return &Store{storeLogger: l, jsonFilePath: path}
}

func (s *Store) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		s.getEntry(responseWriter, request)
		return
	} else if request.Method == http.MethodPost {
		s.addEntry(responseWriter, request)
		return
	} else if request.Method == http.MethodDelete {
		err := repository.Flush(s.jsonFilePath)
		if err != nil {
			http.Error(responseWriter, fmt.Sprintf("An error occurred while flushing the JSON file! Error: %s", err.Error()), http.StatusInternalServerError)
		}
		return
	}

	//HTTP method not supported
	responseWriter.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Store) addEntry(responseWriter http.ResponseWriter, request *http.Request) {
	s.storeLogger.Println("Received HTTP POST request")

	e := repository.Entry{}

	err := e.ConvertFromJSON(request.Body)

	if err != nil {
		http.Error(responseWriter, fmt.Sprintf("An error occurred while processing the data! Error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	repository.AddEntry(e.Key, e.Value, s.storeLogger)
	s.storeLogger.Println("Key value pair added successfully!")

	err = repository.Sync(s.jsonFilePath, s.storeLogger)
	if err != nil {
		s.storeLogger.Println("JSON file could not be synced!")
	}

	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(map[string]string{
		"message": "Key value pair added successfully!",
		"code":    fmt.Sprint(http.StatusOK),
	})
	if err != nil {
		s.storeLogger.Printf("Error happened in JSON marshal. Err: %s", err)
	}
	responseWriter.Write(jsonResp)
}

func (s *Store) getEntry(responseWriter http.ResponseWriter, request *http.Request) {
	s.storeLogger.Println("Received HTTP GET request")

	key := request.URL.Query().Get("key")

	if strings.EqualFold(key, "") {
		http.Error(responseWriter, "An error occurred while processing the data!", http.StatusBadRequest)
		return
	}

	entry, err := repository.GetEntry(key)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	err = entry.ConvertToJSON(responseWriter)
	if err != nil {
		http.Error(responseWriter, "Parse error!", http.StatusInternalServerError)
		return
	}
}
