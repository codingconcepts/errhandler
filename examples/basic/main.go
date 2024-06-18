package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codingconcepts/errhandler"
)

var products = map[string]product{
	"a32fb2bd-b402-4bea-93c2-4a0a567b2261": {
		ID:    "a32fb2bd-b402-4bea-93c2-4a0a567b2261",
		Name:  "a",
		Price: 10.99,
	},
	"b68ed795-0604-4696-8eb2-5b4b927330a0": {
		ID:    "b68ed795-0604-4696-8eb2-5b4b927330a0",
		Name:  "b",
		Price: 20.99,
	},
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("GET /products", errhandler.Wrap(getProducts))
	mux.Handle("GET /products/{id}", errhandler.Wrap(getProduct))

	server := &http.Server{Addr: "localhost:3000", Handler: mux}
	log.Fatal(server.ListenAndServe())
}

func getProducts(w http.ResponseWriter, r *http.Request) error {
	return errhandler.SendJSON(w, products)
}

func getProduct(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	p, ok := products[id]
	if !ok {
		return fmt.Errorf("no product with id: %s", id)
	}

	return errhandler.SendJSON(w, p)
}

type product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
