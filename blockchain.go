package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

type Blockchain struct {
	lastBlockHash []byte
	db            *bolt.DB
}

var (
	blocksBucket = "blocksBucket"
	dbName       = "blockchain.db"
)

func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbName, 0600, nil)
	must(err)
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		if bucket == nil {
			gBlock := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				return err
			}
			bBlock, err := gBlock.Serialize()
			if err != nil {
				return err
			}
			err = b.Put(gBlock.Hash, bBlock)
			if err != nil {
				return err
			}
			err = b.Put([]byte("l"), gBlock.Hash)
			if err != nil {
				return err
			}
			tip = gBlock.Hash
		} else {
			tip = bucket.Get([]byte("l"))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return &Blockchain{
		lastBlockHash: tip,
		db:            db,
	}
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			return errors.New(fmt.Sprintf("blocks bucket %s not exists", blocksBucket))
		}
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Println(err)
	}
	newBlock := NewBlock(data, lastHash)
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		bBlock, err := newBlock.Serialize()
		if err != nil {
			return err
		}
		err = b.Put(newBlock.Hash, bBlock)
		if err != nil {
			return err
		}
		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			return err
		}
		bc.lastBlockHash = newBlock.Hash
		return nil
	})
	
}

func (bc *Blockchain) Iterator() *BlockchainInterator {
	return &BlockchainInterator{
		currentHash: bc.lastBlockHash,
		db:          bc.db,
	}
}
