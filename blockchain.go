package main

import (
	"encoding/hex"
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
)

func NewBlockchain(address, name string) *Blockchain {
	var tip []byte
	db, err := bolt.Open(name, 0600, nil)
	must(err)
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		if bucket == nil {
			coinbaseTx := NewCoinbaseTx(address, "May The Force Be With You")
			gBlock := NewGenesisBlock(coinbaseTx)
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

func NewGenesisBlock(coinbase *Transction) *Block {
	return NewBlock([]*Transction{coinbase}, []byte{})
}

func (bc *Blockchain) AddBlock(txs []*Transction) *Block {
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
		return nil
	}
	newBlock := NewBlock(txs, lastHash)
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
	if err != nil {
		return nil
	}
	return newBlock
}

func (bc *Blockchain) Iterator() *BlockchainInterator {
	return &BlockchainInterator{
		currentHash: bc.lastBlockHash,
		db:          bc.db,
	}
}

func (bc *Blockchain) FindUnspentTransations(address string) []Transction {
	bci := bc.Iterator()
	spentTXs := make(map[string][]int)
	var unspentTXs []Transction
	for {
		block := bci.Next()

		for _, tx := range block.TXs {
			txID := hex.EncodeToString(tx.ID)
		Outputs:
			for outIdx, txOut := range tx.VOut {
				if spentTXs[txID] != nil {
					for _, spentOut := range spentTXs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				if txOut.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}
			if !tx.IsCoinBase() {
				for _, txin := range tx.VIn {
					if txin.CanUnlockOutput(address) {
						inTxID := hex.EncodeToString(txin.Txid)
						spentTXs[inTxID] = append(spentTXs[inTxID], txin.Vout)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return unspentTXs
}

func (bc *Blockchain) FindUTXOs(address string) []TxOutput {
	txs := bc.FindUnspentTransations(address)
	var txOuts []TxOutput

	for _, tx := range txs {
		for _, out := range tx.VOut {
			if out.CanBeUnlockedWith(address) {
				txOuts = append(txOuts, out)
			}
		}
	}

	return txOuts
}

func (bc *Blockchain) FindSpendableUTXOs(address string, amout int) (int, map[string][]int) {
	txOuts := make(map[string][]int)

	txs := bc.FindUnspentTransations(address)

	accumulated := 0

	for _, tx := range txs {
		for outIdx, out := range tx.VOut {
			txID := hex.EncodeToString(tx.ID)
			if out.CanBeUnlockedWith(address) && accumulated < amout {
				accumulated += out.Value
				txOuts[txID] = append(txOuts[txID], outIdx)
				if accumulated >= amout {
					return accumulated, txOuts
				}
			}
		}
	}
	return accumulated, txOuts
}
