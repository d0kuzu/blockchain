package storage

import (
	"blockchain/blockchain"
	"encoding/binary"
	"go.etcd.io/bbolt"
)

type DB struct {
	db *bbolt.DB
}

func InitDB(path string) (*DB, error) {
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func (d *DB) SaveBlock(block *blockchain.Block) error {
	data, err := block.Serialize()
	if err != nil {
		return err
	}

	return d.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("blocks"))
		if err != nil {
			return err
		}
		return bucket.Put(IntToHex(block.Index), data)
	})
}

func (d *DB) LoadBlocks() ([]*blockchain.Block, error) {
	var blocks []*blockchain.Block

	err := d.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("blocks"))
		if bucket == nil {
			return nil
		}
		return bucket.ForEach(func(_, v []byte) error {
			block, err := blockchain.DeserializeBlock(v)
			if err != nil {
				return err
			}
			blocks = append(blocks, block)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	return blocks, nil
}

func IntToHex(num int) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(num))
	return buf
}
