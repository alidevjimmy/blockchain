package main

import (
	"log"

	"github.com/boltdb/bolt"
)

type BlockchainInterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (i *BlockchainInterator) Next() *Block {
	var block *Block
	
	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		rawBlock := b.Get(i.currentHash)
		block = DeserializeBlock(rawBlock)

		return nil
	})
	if err != nil {
		log.Println(err)
		return nil
	}

	i.currentHash = block.PrevBlockHash
	return block
}
