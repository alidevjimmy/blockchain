package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	Hash          []byte
	Version       float64
	Nonce         int

	Data []byte
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Data:          []byte(data),
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
	headers := bytes.Join([][]byte{b.PrevBlockHash, timestamp, version, b.Data}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}
