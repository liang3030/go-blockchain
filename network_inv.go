package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

// It is used to show other nodes what blocks or transactions we have.
// It doesn’t contain whole blocks and transactions, just their hashes.
// Type: whether these are blocks or transactions
type inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

func handleInv(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload inv
	var blocksInTransit [][]byte

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Received inventory with %d items\n", len(payload.Items))

	if payload.Type == "block" {
		blocksInTransit = payload.Items

		blockHash := payload.Items[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}

		for _, b := range blocksInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]

		if mempool[hex.EncodeToString(txID)].ID == nil {
			sendGetData(payload.AddrFrom, "tx", txID)
		}

	}
}

func sendInv(address, kind string, items [][]byte) {
	data := inv{nodeAddress, kind, items}
	payload := gobEncode(data)
	request := append(commandToBytes("inv"), payload...)
	sendData(address, request)
}
