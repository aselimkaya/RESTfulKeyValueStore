package service

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aselimkaya/RESTfulKeyValueStore/src/repository"
)

type Store struct {
	storeLogger  *log.Logger
	jsonFilePath string
}

func New(l *log.Logger, path string) *Store {
	return &Store{storeLogger: l, jsonFilePath: path}
}

func (s *Store) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		//TODO: Add HTTP GET handler
		return
	} else if request.Method == http.MethodPost {
		s.addKeyVal(responseWriter, request)
		return
	} else if request.Method == http.MethodDelete {
		//TODO: Add HTTP DELETE handler
		return
	}

	//HTTP method not supported
	responseWriter.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Store) addKeyVal(responseWriter http.ResponseWriter, request *http.Request) {
	s.storeLogger.Println("Received HTTP POST request")

	e := repository.Entry{}

	err := e.ConvertFromJSON(request.Body)

	if err != nil {
		http.Error(responseWriter, "An error occurred while processing the data!", http.StatusBadRequest)
		return
	}

	if _, ok := repository.KeyValStore[e.Key]; ok {
		s.storeLogger.Println("Key already exists in the store, value will be updated")
	}

	repository.AddPair(e.Key, e.Value)

	jsonString, err := json.Marshal(repository.KeyValStore)
	if err != nil {
		s.storeLogger.Printf("Map updated but JSON file could not be updated. Error: %v\n", err)
		return
	}

	f, err := os.OpenFile(s.jsonFilePath, os.O_TRUNC|os.O_WRONLY, os.FileMode(0666))
	if err != nil {
		s.storeLogger.Fatal(err)
	}
	defer f.Close()

	_, err = f.Write(jsonString)
	if err != nil {
		s.storeLogger.Printf("Map updated but JSON file could not be updated. Error: %v\n", err)
		return
	}

	s.storeLogger.Println("Key value pair added successfully!")

	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(map[string]string{
		"message": "Key value pair added successfully!",
	})
	if err != nil {
		s.storeLogger.Printf("Error happened in JSON marshal. Err: %s", err)
	}
	responseWriter.Write(jsonResp)
}
