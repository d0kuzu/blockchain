package main

import (
	"log"
	"time"

	"blockchain/blockchain"
	"blockchain/network"
	"blockchain/storage"
)

func main() {
	db, err := storage.InitDB("blockchain.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	bc := blockchain.NewBlockchain(db)

	node := network.NewNode("localhost:3000", bc, db)

	go node.StartServer()

	for {
		time.Sleep(10 * time.Second)
		node.CreateBlock()
	}
}
