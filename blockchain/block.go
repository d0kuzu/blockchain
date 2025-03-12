package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"time"
)

type Block struct {
	Index        int
	Timestamp    int64
	Transactions []Transaction
	PrevHash     []byte
	Hash         []byte
	Nonce        int
}

func NewBlock(index int, transactions []Transaction, prevHash []byte) *Block {
	block := &Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		PrevHash:     prevHash,
		Nonce:        0,
	}
	block.Hash = block.CalculateHash()
	return block
}

func (b *Block) CalculateHash() []byte {
	data := bytes.Join([][]byte{
		IntToHex(b.Index),
		IntToHex(int(b.Timestamp)),
		b.PrevHash,
		b.HashTransactions(),
		IntToHex(b.Nonce),
	}, []byte{})

	hash := sha256.Sum256(data)
	return hash[:]
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	for _, tx := range b.Transactions {
		txHash := tx.Hash()
		txHashes = append(txHashes, txHash)
	}
	if len(txHashes) == 0 {
		return []byte{}
	}
	data := bytes.Join(txHashes, []byte{})
	hash := sha256.Sum256(data)
	return hash[:]
}

func (b *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(b); err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}

func DeserializeBlock(data []byte) (*Block, error) {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&block); err != nil {
		return nil, err
	}
	return &block, nil
}

func IntToHex(num int) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(num))
	return buf
}
