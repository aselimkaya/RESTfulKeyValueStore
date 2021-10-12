package service

import (
	"log"
	"net/http"
)

type Store struct {
	storeLogger *log.Logger
}

func New(l *log.Logger) *Store {
	return &Store{l}
}

func (s *Store) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		//TODO: Add HTTP GET handler
		return
	} else if request.Method == http.MethodPost {
		//TODO: Add HTTP POST handler
		return
	} else if request.Method == http.MethodDelete {
		//TODO: Add HTTP DELETE handler
		return
	}

	//HTTP method not supported
	responseWriter.WriteHeader(http.StatusMethodNotAllowed)
}
