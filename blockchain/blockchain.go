package blockchain

import (
	"fmt"

	"github.com/JimmyMcBride/digicoin/utils"
	"github.com/dgraph-io/badger"
)

const (
	dbPath = "/tmp/blocks"
)

// Blockchain is a list of blocks in the blockchain.
type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

// Iterator is struct for the current iteration of the blockchain.
type Iterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// InitBlockchain initializes the blockchain.
func InitBlockchain() *Blockchain {
	var lastHash []byte

	db, err := badger.Open(badger.DefaultOptions("/tmp/blocks"))
	utils.HandleErr(err)

	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			utils.HandleErr(err)
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		}
		item, err := txn.Get([]byte("lh"))
		utils.HandleErr(err)
		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})
		return err
	})
	utils.HandleErr(err)

	blockchain := Blockchain{lastHash, db}
	return &blockchain
}

// AddBlock adds a block to the blockchain.
func (chain *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utils.HandleErr(err)
		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})

		return err
	})
	utils.HandleErr(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		utils.HandleErr(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	utils.HandleErr(err)
}

// Iterator converts a blockchain to a blockchain iterator.
func (chain *Blockchain) Iterator() *Iterator {
	iter := &Iterator{chain.LastHash, chain.Database}

	return iter
}

// Next gives you the "next" block, working backwards through the blockchain.
func (iter *Iterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		utils.HandleErr(err)

		var encodedBlock []byte
		err = item.Value(func(val []byte) error {
			encodedBlock = append([]byte{}, val...)
			return nil
		})
		block = Deserialize(encodedBlock)

		return err
	})

	utils.HandleErr(err)

	iter.CurrentHash = block.PrevHash

	return block
}
