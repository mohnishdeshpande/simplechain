package chain

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/mohnishdeshpande/simplechain/internal/product"
)

type Blockchain struct {
	Blocks []*Block
	mutex  *sync.Mutex
}

var MyBlockchain *Blockchain

func (bc *Blockchain) AddBlock(data product.Checkout) {
	bc.mutex.Lock()

	// ledger will always have atleast one block
	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	// create and validate new block
	newBlock := CreateBlock(prevBlock, data)
	if validBlock(newBlock, prevBlock) {
		bc.Blocks = append(bc.Blocks, newBlock)
	}

	bc.mutex.Unlock()
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: []*Block{CreateGenesisBlock()},
		mutex:  &sync.Mutex{},
	}
}

func GetBlockchain(rw http.ResponseWriter, r *http.Request) {
	payload, err := json.MarshalIndent(MyBlockchain.Blocks, "", " ")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(payload)
		log.Printf("cloud not retrieve blockchain")
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(payload)
}

func WriteBlock(rw http.ResponseWriter, r *http.Request) {
	var data product.Checkout

	// decode the request data
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("could not parse checkout data"))
		log.Printf("could not parse checkout data: %v\n", err)
		return
	}

	// add block to the ledger
	MyBlockchain.AddBlock(data)

	// return new block in json
	payload, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("could not write block"))
		log.Printf("could not write block: %v\n", err)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(payload)
}
