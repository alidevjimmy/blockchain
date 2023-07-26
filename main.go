package main

import (
	"fmt"
	"strconv"
)

func main() {
	bc := NewBlockchain()

	// bc.AddBlock("transfer 1 AC to mmd")
	// bc.AddBlock("transfer 10 AC to jj")

	iter := bc.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Prev Block Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Block Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.IsValid()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
