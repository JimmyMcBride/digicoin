package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

// Block is the struct for each block in the blockchain.
type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

// CreateBlock ...
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// Genesis ...
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

// Serialize ...
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return res.Bytes()
}

// Deserialize ...
func (b *Block) Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(b)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
