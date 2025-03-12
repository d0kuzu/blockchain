package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

type Transaction struct {
	ID     []byte
	From   string
	To     string
	Amount int
}

func (tx *Transaction) Hash() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	_ = encoder.Encode(tx)

	hash := sha256.Sum256(buffer.Bytes())
	return hash[:]
}

func (tx *Transaction) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	_ = encoder.Encode(tx)
	return result.Bytes()
}

func DeserializeTransaction(data []byte) *Transaction {
	var tx Transaction
	decoder := gob.NewDecoder(bytes.NewReader(data))
	_ = decoder.Decode(&tx)
	return &tx
}
