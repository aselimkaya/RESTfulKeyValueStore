package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/aselimkaya/RESTfulKeyValueStore/src/repository"
	"github.com/aselimkaya/RESTfulKeyValueStore/src/service"
)

func main() {
	storeLogger := log.New(os.Stdout, "key-value-store-api", log.LstdFlags)
	jsonFilePath, err := filepath.Abs("./src/db/entries.json")
	if err != nil {
		storeLogger.Fatal(err)
		return
	}

	keyValStoreInit(storeLogger, jsonFilePath)

	handler := service.New(storeLogger, jsonFilePath)

	serveMux := http.NewServeMux()
	serveMux.Handle("/", handler)

	server := &http.Server{
		Addr:        ":80",
		Handler:     serveMux,
		IdleTimeout: 120 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			storeLogger.Fatal(err)
		}
	}()

	storeLogger.Println("Server started successfully at http://localhost:80")

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, syscall.SIGTERM)

	sig := <-signalChannel
	storeLogger.Println("server is shutting down:", sig)

	timeoutContext, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	server.Shutdown(timeoutContext)
}

func keyValStoreInit(l *log.Logger, path string) {
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

		repository.KeyValStore = make(map[string]string)
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

		json.Unmarshal([]byte(byteValue), &repository.KeyValStore)
	}
}
