package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/antonyangani/goservice/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodBye(l)
	ph := handlers.NewProducts(l)

	// creating a new server mux

	sm := http.NewServeMux()
	sm.Handle("/", hh)          // works just like http.HandleFunc which took a http.ResponseWriter and http.Request as parameters
	sm.Handle("/goodbye", gh)   // goodbye handler
	sm.Handle("/products/", ph) // product handler

	// creating a server in go. Its a function that takes some parameters
	// This will enable granular control over timeouts for different types of paylods, idle timeouts, keepalives etc

	server := http.Server{
		Addr:         ":3000",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// creating signal channels
	// these channels will be waiting for kill signals from the os then initiate a shutdown procedure

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)
	// Graceful shutdown: This allows server to finish a transaction before shutdown; this prevents kicking clients off the server amid transaction
	timeoutContex, _ := context.WithTimeout(context.Background(), 30*time.Second) // allows handlers 30s to finish transactions being handled and shutoff the server
	server.Shutdown(timeoutContex)
}
