package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

// MineBlock: mine a block
// It maybe not added to main blockchain that dependes on the length of the blockchain.
func (bc *Blockchain) MineBlock(transactions []*Transaction) *Block {
	var lastHash []byte

	for _, tx := range transactions {
		// verify transaction
		if !bc.VerifyTransaction(tx) {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	err := bc.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte(lastBlockHashKey))

		return nil
	})

	if err != nil {
		fmt.Println("Error in getting a last hash:", err)
	}
	// TODO: check if it is correct
	preBlock := DeserializeBlock(lastHash)
	newBlock := NewBlock(transactions, lastHash, preBlock.Height+1)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			fmt.Println("Error in putting a block:", err)
		}
		err = b.Put([]byte(lastBlockHashKey), newBlock.Hash)
		if err != nil {
			fmt.Println("Error in putting a block:", err)
		}
		bc.tip = newBlock.Hash
		return nil
	})

	if err != nil {
		fmt.Println("Error in updating a db:", err)
	}
	return newBlock
}
