package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// Nodes are communicate with messages.
// When a node runs, it gets several nodes from a DNS seed, and sends them version message.
/*
When a node receives message, tt’ll respond with its own version message.
This is a kind of a handshake: no other interaction is possible without prior greeting of each other.
But it’s not just politeness: version is used to find a longer blockchain.
When a node receives a version message it checks if the node’s blockchain is longer than the value of BestHeight.
If it’s not, the node will request and download missing blocks.
*/

// 1. Version: only have one version
// 2. BestHeight: stores the length of node's blockchain
// 3. AddrFrom: stores the address of sender
type version struct {
	Version    int
	BestHeight int
	AddrFrom   string
}

func handleVersion(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload version

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)

	if err != nil {
		fmt.Println("Error in decoding version:", err)
	}

	mybestHeight := bc.GetBestHeight()
	foreignerBestHeight := payload.BestHeight

	if mybestHeight < foreignerBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if mybestHeight > foreignerBestHeight {
		sendVersion(payload.AddrFrom, bc)
	}

	if !nodeIsKnown(payload.AddrFrom) {
		knownNodes = append(knownNodes, payload.AddrFrom)
	}
}

// messages are sequences of bytes.
// First 12 bytes specify command name (“version” in this case),
// and the latter bytes will contain gob-encoded message structure.
func sendVersion(addr string, bc *Blockchain) {
	bestHeight := bc.GetBestHeight()

	payload := gobEncode(version{nodeVersion, bestHeight, nodeAddress})

	request := append(commandToBytes("version"), payload...)

	sendData(addr, request)
}
