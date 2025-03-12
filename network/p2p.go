package network

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"

	"blockchain/blockchain"
	"blockchain/storage"
)

type Node struct {
	Address    string
	Peers      []string
	Blockchain *blockchain.Blockchain
	DB         *storage.Database
	Mutex      sync.Mutex
}

func NewNode(address string, bc *blockchain.Blockchain, db *storage.Database) *Node {
	return &Node{
		Address:    address,
		Blockchain: bc,
		DB:         db,
	}
}

func (node *Node) StartServer() {
	listener, err := net.Listen("tcp", node.Address)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	fmt.Println("Node is running on", node.Address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection error: %v", err)
			continue
		}
		go node.handleConnection(conn)
	}
}

func (node *Node) handleConnection(conn net.Conn) {
	defer conn.Close()

	var block blockchain.Block
	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&block)
	if err != nil {
		log.Printf("Failed to decode block: %v", err)
		return
	}

	node.Mutex.Lock()
	defer node.Mutex.Unlock()

	node.Blockchain.AddBlock(&block)
	fmt.Printf("New block added: %v\n", block)
}

func (node *Node) CreateBlock() {
	node.Mutex.Lock()
	defer node.Mutex.Unlock()

	newBlock := node.Blockchain.CreateBlock(nil, node.Address)

	for _, peer := range node.Peers {
		go node.broadcastBlock(newBlock, peer)
	}
}

func (node *Node) broadcastBlock(block *blockchain.Block, peer string) {
	conn, err := net.Dial("tcp", peer)
	if err != nil {
		log.Printf("Failed to connect to peer %s: %v", peer, err)
		return
	}
	defer conn.Close()

	encoder := json.NewEncoder(conn)
	err = encoder.Encode(block)
	if err != nil {
		log.Printf("Failed to send block: %v", err)
	}
}
