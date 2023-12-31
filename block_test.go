package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBlock(t *testing.T) {
	block := NewBlock([]*Transction{}, []byte{})
	assert.Equal(t, []*Transction{}, []*Transction{})
	assert.Equal(t, []byte{}, block.PrevBlockHash)
	assert.Equal(t, VERSION, block.Version)
}

func TestSetHash(t *testing.T) {
	prevHash := []byte("prevHash")
	block := NewBlock([]*Transction{}, prevHash)

	block.SetHash()

	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	version := []byte(strconv.FormatFloat(block.Version, 'f', 2, 32))
	headers := bytes.Join([][]byte{block.PrevBlockHash, timestamp, version}, []byte{})
	hash := sha256.Sum256(headers)

	assert.Equal(t, hash[:], block.Hash)
}

func TestSerialize(t *testing.T) {
	prevHash := []byte("prevHash")
	block := NewBlock([]*Transction{}, prevHash)

	s, err := block.Serialize()

	assert.Nil(t, err)
	assert.NotEqual(t, []byte{}, s)
}

func TestDeserializeBlock(t *testing.T) {
	prevHash := []byte("prevHash")
	block := NewBlock([]*Transction{}, prevHash)

	s, err := block.Serialize()
	assert.Nil(t, err)
	assert.NotEqual(t, []byte{}, s)

	dBlock := DeserializeBlock(s)
	assert.Equal(t, block, dBlock)
}
