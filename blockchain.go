package main

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

/*
Blockchain is a database with certain structure.
It is an ordered, back-linked list.
*/

const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "init coinbase data"

// use l as last hash key
const lastBlockHashKey = "l"

// tip is the hash of the last block in the blockchain
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// create a new blockchain DB
// CreateBlockchain:
// 1. create a bucket in database that is used to store blocks of the blockchain.
// 2. create a genesis block and store it in the bucket. It will be the first block in this blockchain.
func CreateBlockchain(address string, nodeID string) *Blockchain {
	dbFile := fmt.Sprintf(dbFile, nodeID)
	if dbExists(dbFile) {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	// TODO: byte type in golang
	var tip []byte

	cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
	genesis := NewGenesisBlock(cbtx)

	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		fmt.Println("Error in open a db file:", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte(lastBlockHashKey), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc

}

// build a blockchain
// NewBlockchain:
// TODO: change the name to build a blockchain if it is necessary
func NewBlockchain(nodeID string) *Blockchain {
	// TODO: check nodeId usage
	fmt.Println("NewBlockchain with nodeId:", nodeID)

	if dbExists(nodeID) == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		fmt.Println("Error in open a db file:", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte(lastBlockHashKey))

		return nil
	})

	if err != nil {
		fmt.Println("Error in updating a db:", err)
	}

	bc := Blockchain{tip, db}
	return &bc
}

func (bc *Blockchain) GetBestHeight() int {
	var lastBlock Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash := b.Get([]byte(lastBlockHashKey))
		blockData := b.Get(lastHash)
		lastBlock = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return lastBlock.Height
}

func dbExists(dbFile string) bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
