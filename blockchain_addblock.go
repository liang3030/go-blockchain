package main

import (
	"log"

	"github.com/boltdb/bolt"
)

// AddBlock: add a block to a blockchain
// Multiple miner are working on mining a block, it uses a longest length of blocks to decide which is the main chain.
func (bc *Blockchain) AddBlock(block *Block) {

	// TODO: previous implementation until chapter 6. Revisit it
	// err = bc.db.Update(func(tx *bolt.Tx) error {
	// 	b := tx.Bucket([]byte(blocksBucket))
	// 	err := b.Put(block.Hash, lastHash)
	// 	if err != nil {
	// 		fmt.Println("Error in putting a block:", err)
	// 	}
	// 	err = b.Put([]byte(lastBlockHashKey), block.Hash)
	// 	if err != nil {
	// 		fmt.Println("Error in putting a block:", err)
	// 	}
	// 	bc.tip = block.Hash
	// 	return nil
	// })

	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		blockInDb := b.Get(block.Hash)

		if blockInDb != nil {
			return nil
		}

		blockData := block.Serialize()
		err := b.Put(block.Hash, blockData)
		if err != nil {
			log.Panic(err)
		}

		lastHash := b.Get([]byte(lastBlockHashKey))
		lastBlockData := b.Get(lastHash)
		lastBlock := DeserializeBlock(lastBlockData)

		if block.Height > lastBlock.Height {
			err = b.Put([]byte(lastBlockHashKey), block.Hash)
			if err != nil {
				log.Panic(err)
			}
			bc.tip = block.Hash
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
