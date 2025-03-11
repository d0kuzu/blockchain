package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Структура блока
type Block struct {
	Index     int
	Timestamp string
	Data      string
	PrevHash  string
	Hash      string
	Validator string // Узел, который создал блок
}

// Структура узла
type Node struct {
	ID    string
	Coins int
}

var (
	blockchain []Block
	nodes      []Node
	mempool    []string
	mutex      sync.Mutex
)

// Функция хеширования блока
func calculateHash(block Block) string {
	record := fmt.Sprintf("%d%s%s%s%s", block.Index, block.Timestamp, block.Data, block.PrevHash, block.Validator)
	h := sha256.New()
	h.Write([]byte(record))
	return hex.EncodeToString(h.Sum(nil))
}

// Выбор узла для создания блока (по минимальному количеству монет)
func selectValidator() *Node {
	mutex.Lock()
	defer mutex.Unlock()

	if len(nodes) == 0 {
		return nil
	}

	// Поиск узлов с минимальным количеством монет
	minCoins := nodes[0].Coins
	var candidates []Node

	for _, node := range nodes {
		if node.Coins < minCoins {
			minCoins = node.Coins
			candidates = []Node{node}
		} else if node.Coins == minCoins {
			candidates = append(candidates, node)
		}
	}

	// Если несколько узлов с одинаковым балансом, выбираем случайный
	selected := candidates[rand.Intn(len(candidates))]
	return &selected
}

// Создание нового блока
func generateBlock(prevBlock Block, validator *Node) Block {
	newBlock := Block{
		Index:     prevBlock.Index + 1,
		Timestamp: time.Now().String(),
		Data:      fmt.Sprintf("Transactions: %v", mempool),
		PrevHash:  prevBlock.Hash,
		Validator: validator.ID,
	}
	newBlock.Hash = calculateHash(newBlock)
	mempool = nil     // Очистка мемпула после создания блока
	validator.Coins++ // Награждаем узел монетами
	return newBlock
}

// Запуск консенсуса
func consensusLoop() {
	for {
		time.Sleep(10 * time.Second)
		validator := selectValidator()
		if validator != nil {
			mutex.Lock()
			newBlock := generateBlock(blockchain[len(blockchain)-1], validator)
			blockchain = append(blockchain, newBlock)
			fmt.Printf("Новый блок #%d создан узлом %s\n", newBlock.Index, newBlock.Validator)
			mutex.Unlock()
		}
	}
}

func main() {
	// Создаем генезис-блок
	genesisBlock := Block{Index: 0, Timestamp: time.Now().String(), Data: "Genesis Block", PrevHash: "", Validator: "System"}
	genesisBlock.Hash = calculateHash(genesisBlock)
	blockchain = append(blockchain, genesisBlock)

	// Добавляем узлы
	nodes = append(nodes, Node{ID: "Node1", Coins: 0})
	nodes = append(nodes, Node{ID: "Node2", Coins: 0})
	nodes = append(nodes, Node{ID: "Node3", Coins: 0})

	// Запускаем процесс консенсуса
	go consensusLoop()

	// Эмуляция поступления транзакций
	for i := 0; i < 5; i++ {
		time.Sleep(3 * time.Second)
		mutex.Lock()
		mempool = append(mempool, fmt.Sprintf("Transaction %d", i+1))
		fmt.Println("Добавлена транзакция", i+1)
		mutex.Unlock()
	}

	// Ожидание для демонстрации работы
	time.Sleep(60 * time.Second)
}
