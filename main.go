package main

import (
	"fmt"
	"log"
	"time"

	"blockchain/blockchain"
	"blockchain/network"
	"blockchain/storage"
)

func main() {
	db, err := storage.InitDB("blockchain.db")
	if err != nil {
		log.Fatal("Ошибка инициализации базы данных:", err)
	}

	bc := blockchain.NewBlockchain(db)

	node := network.NewNode("localhost:5000", db)
	node.Peers = network.NewPeers([]string{"localhost:5001", "localhost:5002"})

	go node.StartServer()

	for {
		time.Sleep(10 * time.Second)

		newBlock, err := bc.AddBlock([]blockchain.Transaction{
			{
				From:   "system",
				To:     "miner",
				Amount: 10,
			},
		})
		if err != nil {
			fmt.Println("Ошибка добавления блока:", err)
			continue
		}

		node.BroadcastBlock(newBlock)

		fmt.Println("Создан и отправлен новый блок:", newBlock.Index)
	}
}
