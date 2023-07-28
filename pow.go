package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

type PoorfOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *PoorfOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-TARGET_BITS))

	pow := &PoorfOfWork{
		block:  b,
		target: target,
	}

	return pow
}

func (pow *PoorfOfWork) prepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.block.HashTransactions(),
		pow.block.PrevBlockHash,
		IntToHex(pow.block.Timestamp),
		IntToHex(int64(TARGET_BITS)),
		IntToHex(int64(nonce)),
	}, []byte{})

	return data
}

// (nonce , hash)
func (pow *PoorfOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	maxNonce := math.MaxInt64

	fmt.Printf("Mining the block which contains \"%d\" transactions\n", len(pow.block.TXs))
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()
	return nonce, hash[:]
}

func (pow *PoorfOfWork) IsValid() bool {
	var hash big.Int
	data := pow.prepareData(pow.block.Nonce)
	rawHash := sha256.Sum256(data)
	hash.SetBytes(rawHash[:])

	return hash.Cmp(pow.target) == -1
}
