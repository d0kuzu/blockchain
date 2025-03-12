package network

import (
	"blockchain/blockchain"
	"blockchain/storage"
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type Node struct {
	Address    string
	Peers      []string
	Blockchain *blockchain.Blockchain
	DB         *storage.DB
	Mutex      sync.Mutex
}

func NewNode(address string, db *storage.DB) *Node {
	bc := blockchain.NewBlockchain(db)
	return &Node{
		Address:    address,
		Blockchain: bc,
		DB:         db,
	}
}

func (n *Node) StartServer() {
	listener, err := net.Listen("tcp", n.Address)
	if err != nil {
		fmt.Println("Ошибка при запуске узла:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Узел запущен на", n.Address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка подключения:", err)
			continue
		}
		go n.HandleConnection(conn)
	}
}

func (n *Node) HandleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	var message Message
	err := decoder.Decode(&message)
	if err != nil {
		fmt.Println("Ошибка декодирования сообщения:", err)
		return
	}

	n.HandleMessage(message)
}

func (n *Node) HandleMessage(message Message) {
	switch message.Type {
	case "new_block":
		var block blockchain.Block
		err := json.Unmarshal(message.Data, &block)
		if err != nil {
			fmt.Println("Ошибка разбора блока:", err)
			return
		}

		n.Mutex.Lock()
		prevBlock := n.Blockchain.GetLastBlock()
		if n.Blockchain.IsValidNewBlock(&block, prevBlock) {
			_, err := n.Blockchain.AddBlock(block.Transactions)
			if err != nil {
				return
			}
			fmt.Println("Блок добавлен в цепочку:", block.Index)
		}
		n.Mutex.Unlock()

	case "new_peer":
		var newPeer string
		err := json.Unmarshal(message.Data, &newPeer)
		if err != nil {
			fmt.Println("Ошибка разбора узла:", err)
			return
		}
		n.AddPeer(newPeer)
	}
}

func (n *Node) AddPeer(peer string) {
	for _, existingPeer := range n.Peers {
		if existingPeer == peer {
			return
		}
	}
	n.Peers = append(n.Peers, peer)
	fmt.Println("Добавлен новый узел:", peer)
}

func (n *Node) Broadcast(message Message) {
	for _, peer := range n.Peers {
		conn, err := net.Dial("tcp", peer)
		if err != nil {
			fmt.Println("Ошибка соединения с узлом:", peer)
			continue
		}
		defer conn.Close()

		encoder := json.NewEncoder(conn)
		err = encoder.Encode(message)
		if err != nil {
			fmt.Println("Ошибка отправки сообщения:", err)
		}
	}
}

func (n *Node) BroadcastBlock(block *blockchain.Block) {
	data, _ := json.Marshal(block)
	message := Message{
		Type: "new_block",
		Data: data,
	}
	n.Broadcast(message)
}
