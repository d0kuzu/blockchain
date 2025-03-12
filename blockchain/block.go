package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"blockchain/storage"
)

type Block struct {
	Index        int
	Timestamp    int64
	Transactions []Transaction
	PrevHash     string
	Hash         string
	Validator    string
}

type Transaction struct {
	Sender    string
	Recipient string
	Amount    float64
}

type Blockchain struct {
	Blocks []*Block
	DB     *storage.Database
}

func NewBlockchain(db *storage.Database) *Blockchain {
	bc := &Blockchain{DB: db}
	genesis := bc.CreateGenesisBlock()
	bc.AddBlock(genesis)
	return bc
}

func (bc *Blockchain) CreateGenesisBlock() *Block {
	return &Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Transactions: []Transaction{},
		PrevHash:     "",
		Hash:         calculateHash(0, time.Now().Unix(), nil, ""),
		Validator:    "system",
	}
}

func (bc *Blockchain) AddBlock(block *Block) {
	bc.Blocks = append(bc.Blocks, block)
	err := bc.DB.SaveBlock(block)
	if err != nil {
		log.Printf("Failed to save block: %v", err)
	}
}

func (bc *Blockchain) CreateBlock(transactions []Transaction, validator string) *Block {
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := &Block{
		Index:        lastBlock.Index + 1,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		PrevHash:     lastBlock.Hash,
		Hash:         calculateHash(lastBlock.Index+1, time.Now().Unix(), transactions, lastBlock.Hash),
		Validator:    validator,
	}
	bc.AddBlock(newBlock)
	return newBlock
}

func calculateHash(index int, timestamp int64, transactions []Transaction, prevHash string) string {
	record := fmt.Sprintf("%d%d%s%s", index, timestamp, transactionsToString(transactions), prevHash)
	h := sha256.New()
	h.Write([]byte(record))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func transactionsToString(transactions []Transaction) string {
	data, _ := json.Marshal(transactions)
	return string(data)
}
