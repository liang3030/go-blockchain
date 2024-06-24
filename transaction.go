package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy = 10

/*
Transaction:
1. There are outputs that are not linked to inputs.
2. In one transaction, inputs can reference ouptuts from multiple transactions.
3. An input must reference an output.
*/

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

/*
TXInput:references a previous output.
1. Txid stores previous output transaction transaction ID. That transactions contains the output.
2. Vout is an index of the previous output in the transaction.
3. PubKey is a public key.
*/
type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}

/*
TXOutput:
1. Value is the amount of coins.
2. PubKeyHash is a hashed public key.

outputs is not divisible. You cannot reference a part of its value.
*/
type TXOutput struct {
	Value      int
	PubKeyHash []byte
}

/*
When start to mint a block, it adds a coinbase transaction to it.
Coinbase transaction is a special transaction does not need previous existing output.
*/

func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	// Coinbase transaction has only one input.
	// Txid is empty and Vout is -1, and does not store ScriptSig.
	txin := TXInput{[]byte{}, -1, nil, []byte(data)}

	// subsidy is the amout of reward.
	txout := TXOutput{subsidy, []byte(to)}

	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()
	return &tx
}

/*
Send some coins to someone:
1. Create a new transaction
2. Put it in a block
3. Mine the block
*/
func NewUTXOTransaction(wallet *Wallet, to string, amount int, UTXOSet *UTXOSet) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	pubKeyHash := HashPubKey(wallet.PublicKey)
	acc, validOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("Error: Not enough funds")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)

		if err != nil {
			log.Panic(err, "Error in NewUTXOTransaction")
		}

		for _, out := range outs {
			input := TXInput{txID, out, nil, nil}
			inputs = append(inputs, input)
		}
	}

	from := fmt.Sprintf("%s", wallet.GetAddress())
	outputs = append(outputs, *NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc-amount, from)) // a change
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	UTXOSet.Blockchain.SignTransaction(&tx, wallet.PrivateKey)
	return &tx
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		fmt.Println("Error in SetID:", err)
		log.Panic(err)
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// check input use a specific key to unlock the output.
// Notice: inputs store raw public keys, but the function takes hashed public keys.
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {

	lockingHash := HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

// check if provided public key hash was used to lock the output.
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Equal(out.PubKeyHash, pubKeyHash)
}

// Lock simply locks an output.
// When we send coins to someone, we know only their address, thus the function takes an address as the only argument.
// The address is then decoded and the public key hash is extracted from it and saved in the PubKeyHash field.
// address: The transaction senders's address.
// TODO: where to use this function?
func (out *TXOutput) Lock(address []byte) {
	pubKeyHash, err := Base58Decode(address)
	if err != nil {
		log.Panic("Error in Lock")
	}
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

func DeserializeTransaction(data []byte) Transaction {
	var transaction Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}

	return transaction
}
