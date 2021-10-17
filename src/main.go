package main

import (
	"context"
	"fmt"
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

	//Looking up to environment to find out PORT parameter
	port, found := os.LookupEnv("PORT")
	if !found {
		port = "8000"
	}

	server := &http.Server{
		Addr:        fmt.Sprintf(":%v", port),
		Handler:     serveMux,
		IdleTimeout: 120 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			storeLogger.Fatal(err)
		}
	}()

	storeLogger.Printf("Server started successfully at http://localhost:%v", port)

	//Interrupt or kill signal handler
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, syscall.SIGTERM)

	sig := <-signalChannel
	storeLogger.Println("server is shutting down:", sig)

	timeoutContext, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()
	server.Shutdown(timeoutContext)
}
