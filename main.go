package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("Hello world")
		d, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, "Oops", http.StatusBadGateway)
			return
		}
		fmt.Fprintf(rw, "Hello %s", d)
		// log.Printf("The data: %s\n", d)
	})
	http.HandleFunc("/goodbye", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("Goodbye world")
	})
	http.ListenAndServe(":3000", nil)
}
