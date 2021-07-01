package handlers

import (
	"log"
	"net/http"

	"github.com/antonyangani/goservice/data"
)

type Products struct {
	l *log.Logger
}

// initializes the object Products capable of doing logs

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// GET
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	// PUT (UPDATE)

	if r.Method == http.MethodPut {

	}
	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	listOfProducts := data.GetProducts()
	// converting listOfProducts into JSON
	err := listOfProducts.ToJSON(rw)
	if err != nil {
		http.Error(rw, "unable to marshall json", http.StatusInternalServerError)
	}
}
