package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type addr struct {
	AddrList []string
}

func handleAddr(request []byte) {
	var buff bytes.Buffer
	var payload addr

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	knownNodes = append(knownNodes, payload.AddrList...)
	fmt.Printf("There are %d known nodes now!\n", len(knownNodes))
	requestBlocks()
}

func requestBlocks() {
	for _, node := range knownNodes {
		sendGetBlocks(node)
	}
}
