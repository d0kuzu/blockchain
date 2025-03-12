package storage

import (
	"blockchain/blockchain"
	"encoding/json"
	bolt "go.etcd.io/bbolt"
)

type Database struct {
	db *bolt.DB
}

func InitDB(dbPath string) (*Database, error) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) SaveBlock(block *blockchain.Block) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("blocks"))
		if err != nil {
			return err
		}

		blockData, err := json.Marshal(block)
		if err != nil {
			return err
		}

		return b.Put([]byte(block.Hash), blockData)
	})
}

func (d *Database) LoadBlockchain() ([]*blockchain.Block, error) {
	var blocks []*blockchain.Block

	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("blocks"))
		if b == nil {
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			var block blockchain.Block
			if err := json.Unmarshal(v, &block); err != nil {
				return err
			}
			blocks = append(blocks, &block)
			return nil
		})
	})

	return blocks, err
}
