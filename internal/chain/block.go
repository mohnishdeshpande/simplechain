package chain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mohnishdeshpande/simplechain/internal/product"
)

type Block struct {
	Pos       int
	Data      product.Checkout
	Timestamp string
	Hash      string
	PrevHash  string
}

func (b *Block) generateHash() string {
	// combine values to compute hash
	bytes, _ := json.Marshal(b.Data)
	hashData := string(rune(b.Pos)) + string(bytes) + b.Timestamp + b.PrevHash

	// create new hash
	hash := sha256.New()
	hash.Write([]byte(hashData))

	// write to field
	return hex.EncodeToString(hash.Sum(nil))
}

func (b *Block) validateHash() bool {
	return b.Hash == b.generateHash()
}

func (b *Block) Print() {
	fmt.Printf("Pos: %d\n", b.Pos)
	fmt.Printf("Prev Hash: %x\n", b.PrevHash)
	fmt.Printf("Hash: %x\n", b.Timestamp)
	bytes, _ := json.MarshalIndent(b.Data, "", " ")
	fmt.Printf("Data: %v\n", string(bytes))
	fmt.Printf("Hash: %x\n", b.Hash)
}

// Non struct functions
func CreateGenesisBlock() *Block {
	return CreateBlock(&Block{}, product.Checkout{IsGenesis: true})
}

func CreateBlock(prevBlock *Block, data product.Checkout) *Block {
	block := &Block{
		Pos:       prevBlock.Pos + 1,
		Data:      data,
		Timestamp: time.Now().String(),
		PrevHash:  prevBlock.Hash,
	}
	block.Hash = block.generateHash()

	return block
}

func validBlock(block *Block, prevBlock *Block) bool {
	return block.Pos == prevBlock.Pos+1 && block.PrevHash == prevBlock.Hash && block.validateHash()
}
