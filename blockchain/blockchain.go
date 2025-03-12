package blockchain

import (
	"blockchain/storage"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Blockchain struct {
	Blocks []*Block
	DB     *storage.DB
}

func NewBlockchain(db *storage.DB) *Blockchain {
	bc := &Blockchain{DB: db}
	bc.LoadBlockchain()
	return bc
}

func (bc *Blockchain) AddBlock(transactions []Transaction) (*Block, error) {
	if len(bc.Blocks) == 0 {
		return nil, fmt.Errorf("нет предыдущего блока")
	}

	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(prevBlock.Index+1, transactions, prevBlock.Hash)

	if !bc.IsValidNewBlock(newBlock, prevBlock) {
		return nil, fmt.Errorf("блок недействителен")
	}

	bc.Blocks = append(bc.Blocks, newBlock)

	err := bc.DB.SaveBlock(newBlock)
	if err != nil {
		return nil, err
	}

	return newBlock, nil
}

func (bc *Blockchain) LoadBlockchain() error {
	blocks, err := bc.DB.LoadBlocks()
	if err != nil {
		return err
	}
	bc.Blocks = blocks
	return nil
}

func (bc *Blockchain) GetLastBlock() *Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

func (bc *Blockchain) IsValidNewBlock(newBlock, prevBlock *Block) bool {
	if !bytes.Equal(newBlock.PrevHash, prevBlock.Hash) {
		fmt.Println("Ошибка: хеш предыдущего блока не совпадает")
		return false
	}

	calculatedHash := newBlock.CalculateHash()
	if !bytes.Equal(newBlock.Hash, calculatedHash) {
		fmt.Println("Ошибка: некорректный хеш блока")
		return false
	}

	return true
}

func calculateBlockHash(block *Block) string {
	data := fmt.Sprintf("%d%s%d", block.Index, block.PrevHash, len(block.Transactions))
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
