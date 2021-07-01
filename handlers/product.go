package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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

	// POST (Adding a new Product)

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	// PUT (Updating product information)

	if r.Method == http.MethodPut {

		p.l.Println("Handling PUT requests")
		// expects to get the ID in the URI

		// we will be use a regular expression to get the ID of the product
		rg := regexp.MustCompile(`/([0-9]+)`)
		gp := rg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(gp) != 1 {
			p.l.Println("Got a wrong URL", gp)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		if len(gp[0]) != 2 {
			p.l.Println("Got more than one group", gp)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		idString := gp[0][1]
		// convert into an integer

		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("unable to convert", idString)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		p.updateProduct(id, rw, r)

	}

	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handles GET requests")
	listOfProducts := data.GetProducts()
	// converting listOfProducts into JSON
	err := listOfProducts.ToJSON(rw)
	if err != nil {
		http.Error(rw, "unable to marshall json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handles POST for new products")
	newProduct := &data.Product{}
	err := newProduct.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Could not unmarshall json", http.StatusBadRequest)
	}
	p.l.Printf("Prod: %#v", newProduct)
	data.AddProduct(newProduct)
}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handles PUT for products")
	newProduct := &data.Product{}
	err := newProduct.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Could not unmarshall json", http.StatusBadRequest)
	}
	p.l.Printf("Prod: %#v", newProduct)
	er := data.UpdateProduct(id, newProduct)
	if er == data.ErrorProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	} else if er != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
