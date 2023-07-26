package main

import (
	"fmt"
	"strconv"
)

func main() {
	bc := NewBlockchain()

	bc.AddBlock("transfer 1 AT to mmd")
	bc.AddBlock("transfer 10 AT to jj")

	for _, block := range bc.Blocks {
		fmt.Printf("Prev Block Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Block Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.IsValid()))
		fmt.Println()
	}
}
