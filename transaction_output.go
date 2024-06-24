package main

import (
	"bytes"
	"encoding/gob"
)

/*
UTXO: unspent transaction output. These outputs were not referenced in any inputs.
*/

// TXOutputs collects TXOutput
type TXOutputs struct {
	Outputs []TXOutput
}

// NewTXOutput create a new TXOutput
func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

func (outs TXOutputs) Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		panic(err)
	}

	return buff.Bytes()
}
