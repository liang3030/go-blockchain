package main

import (
	"fmt"

	"github.com/boltdb/bolt"
)

// TODO: implement GetBlock
// TODO: check if it is correct - current it is automatically generated
func (bc *Blockchain) GetBlock(hash []byte) (Block, error) {
	var block Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(hash)
		block = *DeserializeBlock(encodedBlock)
		return nil
	})

	if err != nil {
		fmt.Println("Error in getting a block:", err)
	}

	return block, nil
}

func (bc *Blockchain) GetBlockHashes() [][]byte {
	var blocks [][]byte
	bci := bc.Iterator()

	for {
		block := bci.Next()

		blocks = append(blocks, block.Hash)

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return blocks
}
