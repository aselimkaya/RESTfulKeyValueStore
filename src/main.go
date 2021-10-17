package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aselimkaya/RESTfulKeyValueStore/src/service"
)

func main() {
	storeLogger := log.New(os.Stdout, "key-value-store-api", log.LstdFlags)
	path, err := os.Getwd()
	if err != nil {
		storeLogger.Fatal(err)
	}

	handler := service.New(storeLogger, path)

	serveMux := http.NewServeMux()
	serveMux.Handle("/", handler)
	serveMux.Handle("/entry", handler)

	server := &http.Server{
		Addr:        ":8000",
		Handler:     serveMux,
		IdleTimeout: 120 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			storeLogger.Fatal(err)
		}
	}()

	storeLogger.Println("Server started successfully at http://localhost:8000")

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, syscall.SIGTERM)

	sig := <-signalChannel
	storeLogger.Println("server is shutting down:", sig)

	timeoutContext, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	server.Shutdown(timeoutContext)
}
