package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mohnishdeshpande/simplechain/internal/chain"
	"github.com/mohnishdeshpande/simplechain/internal/product"
)

func main() {
	chain.MyBlockchain = chain.NewBlockchain()

	// assign routes
	router := mux.NewRouter()
	router.HandleFunc("/", chain.GetBlockchain).Methods("GET")
	router.HandleFunc("/", chain.WriteBlock).Methods("POST")
	router.HandleFunc("/new", product.NewProduct).Methods("POST")

	// print the whole blockchain (goroutine)
	go func() {
		for _, block := range chain.MyBlockchain.Blocks {
			block.Print()
		}
	}()

	log.Println("Listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
