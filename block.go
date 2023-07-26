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

	Data []byte
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	return &Block{
		Data:          []byte(data),
		Version:       VERSION,
		PrevBlockHash: prevBlockHash,
		Timestamp:     time.Now().Unix(),
		Hash:          []byte{},
	}
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	version := []byte(strconv.FormatFloat(b.Version, 'f', 2, 32))
	headers := bytes.Join([][]byte{b.PrevBlockHash, timestamp, version, b.Data}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}
