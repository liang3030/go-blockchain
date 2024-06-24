package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type block struct {
	AddrFrom string
	Block    []byte
}

// it means show me what blocks you have
type getBlocks struct {
	AddrForm string
}

func handleBlock(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload block

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := DeserializeBlock(blockData)

	fmt.Println("Received a new block!")
	// TODO: update this function
	bc.AddBlock(block)

	fmt.Printf("Added block %x\n", block.Hash)

	// TODO: check this function later
	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := UTXOSet{bc}
		UTXOSet.Reindex()
	}

}

func handleGetBlocks(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload getBlocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()
	sendInv(payload.AddrForm, "block", blocks)
}

func sendGetBlocks(addr string) {
	payload := gobEncode(getBlocks{nodeAddress})
	request := append(commandToBytes("getblocks"), payload...)
	sendData(addr, request)
}

func sendBlock(addr string, b *Block) {
	data := block{nodeAddress, b.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("block"), payload...)
	sendData(addr, request)
}
