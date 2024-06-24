package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

/*
blockchain network is a P2P(peer-to-peer) network.
It means that nodes are connected direcly to each other.
It is topology is flat, since there are no hirarchy in node roles.
*/

/*
Node Roles
1. Miner
Such nodes are run on powerful or specialized hardware(like ASIC),and their only goal is to mine new blocks as fast as possible. Miners are only possible in blockchains that use Proof-of-Work, because mining actually meanssolving PoW puzzles. In Proof-of-Stake blockchains, there are no miners, but validators instead.


2. Full Node
These nodes validate blocks mined by miners and verify transactions. To do this, they must have the whole copy of blockchain. Also, such nodes perform such routing operations, like helping other nodes to discover each other.

It is very crucial for network to have many full nodes, because it is these nodes that make decisions: they decide if a block or transaction is valid.


3. SPV
SPV stands for Simplified Payment Verification. These nodes do not store a full copy of blockchain, but they still able to verify transactions(not all of them, but a subset). An SPV node depends on a full node to get data from, and there could be many SPV nodes connected to one full node. SPV makes wallet applications possible: one don’t need to download full blockchain, but still can verify their transactions.

*/

/*
The Scenario

1. The central node creates a blockchain.
2. Other (wallet) node connects to it and downloads the blockchain.
3. One more (miner) node connects to the central node and downloads the blockchain.
4. The wallet node creates a transaction.
5. The miner nodes receives the transaction and keeps it in its memory pool.
6. When there are enough transactions in the memory pool, the miner starts mining a new block.
7. When a new block is mined, it’s send to the central node.
9. The wallet node synchronizes with the central node.
9. User of the wallet node checks that their payment was successful.
10. This is what it looks like in Bitcoin. Even though we’re not going to build a real P2P network, we’re going to implement a real, and the main and most important, use case of Bitcoin.
*/

// TODO: block, tx, data, inv -> relationship
var protocol = "tcp"
var nodeVersion = 1

const commandLength = 12

var mempool = make(map[string]Transaction)

var blocksInTransit = [][]byte{}

var minerAddress string

var nodeAddress string
var knownNodes = []string{"localhost:3000"}

func StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	ln, err := net.Listen(protocol, nodeAddress)

	if err != nil {
		fmt.Println("Error in listening:", err)
	}

	defer ln.Close()

	bc := NewBlockchain(nodeID)

	// if current node is not the central one,
	// it must send version message to the central node to find out if its blockchain is outdated.
	if nodeAddress != knownNodes[0] {
		sendVersion(knownNodes[0], bc)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println("Error in accepting a connection:", err)
		}

		go handleConnection(conn, bc)
	}

}

func handleConnection(conn net.Conn, bc *Blockchain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
	case "addr":
		handleAddr(request)
	case "block":
		handleBlock(request, bc)
	case "inv":
		handleInv(request, bc)
	case "getblocks":
		handleGetBlocks(request, bc)
	case "getdata":
		handleGetData(request, bc)
	case "tx":
		handleTx(request, bc)
	case "version":
		handleVersion(request, bc)
	default:
		fmt.Println("Unknown command!")
	}

	conn.Close()
}

// TODO: implement this function
func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// message on the lower level are sequences of bytes.
// First 12 bytes specify command name, and the latter bytes will contain gob-encoded message structure.
func commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}
	return bytes[:]
}

// When a node receives a command, it runs bytesToCommand function to extract the command name and processes command body with correct handler.
func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return string(command)
}

func nodeIsKnown(addr string) bool {
	for _, node := range knownNodes {
		if node == addr {
			return true
		}
	}
	return false
}
