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
	handler := service.New(storeLogger)

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

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, syscall.SIGTERM)

	sig := <-signalChannel
	storeLogger.Println("server is shutting down:", sig)

	timeoutContext, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	server.Shutdown(timeoutContext)
}
