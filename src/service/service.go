package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aselimkaya/RESTfulKeyValueStore/src/repository"
)

//Store holds server side logger and DB JSON file's path inside.
type Store struct {
	storeLogger  *log.Logger
	jsonFilePath string
}

//New function first creates required files and folders and returns a new Store object.
func New(l *log.Logger, path string) *Store {
	os.Mkdir(path+"/db", 0755)
	repository.Init(l, path+"/db/entries.json")
	return &Store{storeLogger: l, jsonFilePath: path + "/db/entries.json"}
}

//ServeHTTP is main handler that handles every single HTTP request. It supports only GET request for welcome page but GET, POST and DELETE for /entry page
func (s *Store) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if strings.EqualFold(request.URL.Path, "/") {
		if request.Method == http.MethodGet {
			s.setResponse(responseWriter, http.StatusOK, "Welcome!")
			return
		}
		//HTTP method not supported
		s.setResponse(responseWriter, http.StatusBadRequest, "HTTP method not supported!")
		return
	} else if strings.EqualFold(request.URL.Path, "/entry") {
		if request.Method == http.MethodGet {
			s.GetEntry(responseWriter, request)
			return
		} else if request.Method == http.MethodPost {
			s.AddEntry(responseWriter, request)
			return
		} else if request.Method == http.MethodDelete {
			err := repository.Flush(s.jsonFilePath)
			if err != nil {
				s.setResponse(responseWriter, http.StatusInternalServerError, fmt.Sprintf("An error occurred while flushing the JSON file! Error: %s", err.Error()))
				return
			}
			s.setResponse(responseWriter, http.StatusOK, "JSON file flushed successfully!")
			return
		}

		//HTTP method not supported
		s.setResponse(responseWriter, http.StatusBadRequest, "HTTP method not supported!")
	}

	s.setResponse(responseWriter, http.StatusNotFound, "Page not found!")
}

//AddEntry is a service function that extracts fields from request body and sends them to repository function
func (s *Store) AddEntry(responseWriter http.ResponseWriter, request *http.Request) {
	s.storeLogger.Println("Received HTTP POST request")

	e := repository.Entry{}

	err := e.ConvertFromJSON(request.Body)

	if err != nil {
		s.setResponse(responseWriter, http.StatusBadRequest, fmt.Sprintf("An error occurred while processing the data! Error: %s", err.Error()))
		return
	}

	if len(e.Key) == 0 {
		s.setResponse(responseWriter, http.StatusBadRequest, `Missing field! 'key' field is required!`)
		return
	} else if len(e.Value) == 0 {
		s.setResponse(responseWriter, http.StatusBadRequest, `Missing field! 'value' field is required!`)
		return
	}

	isExists := repository.AddEntry(e.Key, e.Value, s.storeLogger)
	if isExists {
		s.setResponse(responseWriter, http.StatusOK, "Key already exists, value will be upated")
	} else {
		s.setResponse(responseWriter, http.StatusOK, "Key value pair added successfully")
		s.storeLogger.Println("Key value pair added successfully!")
	}

	err = repository.Sync(s.jsonFilePath, s.storeLogger)
	if err != nil {
		s.storeLogger.Println("JSON file could not be synced!")
	}
}

//GetEntry is a service function that extracts the key from request parameters and sends it to repository to find if the key exists
func (s *Store) GetEntry(responseWriter http.ResponseWriter, request *http.Request) {
	s.storeLogger.Println("Received HTTP GET request")

	key := request.URL.Query().Get("key")

	if strings.EqualFold(key, "") {
		s.setResponse(responseWriter, http.StatusBadRequest, "An error occurred while processing the data!")
		return
	}

	entry, err := repository.GetEntry(key)
	if err != nil {
		s.setResponse(responseWriter, http.StatusBadRequest, fmt.Sprintf("An error occurred while processing the data! Error: %s", err.Error()))
		return
	}

	b, _ := json.Marshal(entry)

	s.setResponse(responseWriter, http.StatusOK, string(b))
}

//setResponse is a utility function thet generates HTTP Responses
func (s *Store) setResponse(responseWriter http.ResponseWriter, status int, message string) http.ResponseWriter {
	responseWriter.WriteHeader(status)
	responseWriter.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(map[string]interface{}{
		"message": message,
		"status":  status,
	})
	if err != nil {
		s.storeLogger.Printf("Error happened in JSON marshal. Err: %s", err)
	}
	responseWriter.Write(jsonResp)

	return responseWriter
}
