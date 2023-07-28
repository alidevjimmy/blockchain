package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	Hash          []byte
	Version       float64
	Nonce         int

	TXs []*Transction
}

func NewBlock(txs []*Transction, prevBlockHash []byte) *Block {
	block := &Block{
		TXs:           txs,
		Version:       VERSION,
		PrevBlockHash: prevBlockHash,
		Timestamp:     time.Now().Unix(),
		Hash:          []byte{},
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return block
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	version := []byte(strconv.FormatFloat(b.Version, 'f', 2, 32))
	headers := bytes.Join([][]byte{b.PrevBlockHash, timestamp, version}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

func (b *Block) Serialize() ([]byte, error) {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	if err := encoder.Encode(b); err != nil {
		return []byte{}, err
	}

	return res.Bytes(), nil
}

func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	if err := decoder.Decode(&block); err != nil {
		return nil
	}
	return &block
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.TXs {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}
