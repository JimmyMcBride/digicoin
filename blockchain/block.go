package blockchain

import (
	"bytes"
	"encoding/gob"

	"github.com/JimmyMcBride/digicoin/utils"
)

// Block is the struct for each block in the blockchain.
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

// CreateBlock creates a new block in the blockchain.
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// Genesis creates the blockchains first block.
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

// Serialize encodes a block.
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	utils.HandleErr(err)

	return res.Bytes()
}

// Deserialize takes in encoded data and returns an un-encoded block.
func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	utils.HandleErr(err)

	return &block
}
