package main

import (
	"fmt"

	"github.com/boltdb/bolt"
)

/*
An iterator is created each time when we want to iterate over blocks in a blockchain.
It will store the block hash of the current iteration and a connection to a DB
An iterator is attached to a blockchain. It is a Blockchain instance that stores a DB connection.
*/

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Iterator returns a BlockchainIterator
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

// Iterator iterate a blockchain form last to first block.
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	// Provide a function that will be executed in a read-only mode.
	// Create a read-only object / instance tx
	err := i.db.View(func(tx *bolt.Tx) error {
		// Get a bucket from the tx object
		b := tx.Bucket([]byte(blocksBucket))

		// Read a key-value pair from the bucket
		encodedBlock := b.Get(i.currentHash)

		block = DeserializeBlock(encodedBlock)
		return nil
	})

	if err != nil {
		fmt.Println("Error in getting a block:", err)
	}

	i.currentHash = block.PrevBlockHash
	return block
}
