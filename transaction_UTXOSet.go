package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"log"

	"github.com/boltdb/bolt"
)

/*
Having the UTXO set means that our data (transactions) are now split into to storages:
1. actual transactions are stored in the blockchain,
2. unspent outputs are stored in the UTXO set.
Such separation requires solid synchronization mechanism
because we want the UTXO set to always be updated and store outputs of most recent transactions.
*/

const utxoBucket = "chainstate"

// unspent transaction outputs: UTXO
// a set of unspent transaction outputs.
type UTXOSet struct {
	// TODO: why this struct store a blockchain pointer?
	Blockchain *Blockchain
}

// Create initial UTXO set
func (u UTXOSet) Reindex() {
	db := u.Blockchain.db
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		// remove the bucket if it exists
		err := tx.DeleteBucket(bucketName)
		if err != nil && err != bolt.ErrBucketNotFound {
			return err
		}
		_, err = tx.CreateBucket(bucketName)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	// get all unspent outputs from blockchain
	UTXO := u.Blockchain.FindUTXO()

	// save the outputs to the bucket
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				return err
			}
			err = b.Put(key, outs.Serialize())
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

// Used for sending coins or check balance
func (u UTXOSet) FindSpendableOutputs(pubkeyHas []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := DeserializeOutputs(v)

			for outIdx, out := range outs.Outputs {
				if out.IsLockedWithKey(pubkeyHas) && accumulated < amount {
					accumulated += out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)

	}
	return accumulated, unspentOutputs
}

// Find unspent transaction outputs by provided public key hash.
func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	db := u.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := DeserializeOutputs(v)

			for _, out := range outs.Outputs {
				if pubKeyHash == nil || out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}

		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return UTXOs
}

// But we don’t want to reindex every time a new block is mined because it’s these frequent blockchain scans that we want to avoid.
// Thus, we need a mechanism of updating the UTXO set:
func (u UTXOSet) Update(block *Block) {
	db := u.Blockchain.db
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		for _, tx := range block.Transactions {
			if tx.IsCoinbase() == false {
				for _, vin := range tx.Vin {
					updatedOuts := TXOutputs{}
					outsBytes := b.Get(vin.Txid)
					outs := DeserializeOutputs(outsBytes)

					for outIdx, out := range outs.Outputs {
						if outIdx != vin.Vout {
							updatedOuts.Outputs = append(updatedOuts.Outputs, out)
						}
					}

					if len(updatedOuts.Outputs) == 0 {
						err := b.Delete(vin.Txid)
						if err != nil {
							return err
						}
					} else {
						err := b.Put(vin.Txid, updatedOuts.Serialize())
						if err != nil {
							return err
						}
					}
				}
			}

			newOutputs := TXOutputs{}
			for _, out := range tx.Vout {
				newOutputs.Outputs = append(newOutputs.Outputs, out)
			}

			err := b.Put(tx.ID, newOutputs.Serialize())
			if err != nil {
				return err
			}

		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}

func (u UTXOSet) CountTransactions() int {
	db := u.Blockchain.db
	counter := 0

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			counter++
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return counter
}

// TODO: check this function again
func DeserializeOutputs(v []byte) TXOutputs {
	var outputs TXOutputs
	decoder := gob.NewDecoder(bytes.NewReader(v))
	err := decoder.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}
	return outputs
}
