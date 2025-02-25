package product

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Product struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Seller     string `json:"seller"`
	UploadDate string `json:"upload_date"`
	SerialNum  string `json:"serial_num"`
}

type Checkout struct {
	ProductId    string `json:"product_id"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis    bool   `json:"is_genesis"`
}

func NewProduct(rw http.ResponseWriter, r *http.Request) {
	var product Product

	// cant decode error
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("could not create a new product"))
		log.Printf("could not create: %v\n", err)
		return
	}

	// assign proudct id - computing hash
	h := md5.New()
	io.WriteString(h, product.SerialNum+product.UploadDate)
	product.Id = fmt.Sprintf("%x", h.Sum(nil))

	// convert Product struct to JSON
	payload, err := json.MarshalIndent(product, "", " ")

	// cant marshal error
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("could not save product data"))
		log.Printf("could not marshal payload: %v\n", err)
		return
	}

	// finally send payload
	rw.WriteHeader(http.StatusOK)
	rw.Write(payload)
}
