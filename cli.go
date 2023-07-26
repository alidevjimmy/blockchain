package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	bc *Blockchain
}

func NewCLI(bc *Blockchain) *CLI {
	return &CLI{
		bc: bc,
	}
}

func (cli *CLI) Run() error {
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printchainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addblockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			return err
		}
	case "printchain":
		err := printchainCmd.Parse(os.Args[2:])
		if err != nil {
			return err
		}
	default:
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addblockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addblockData)
	}
	if printchainCmd.Parsed() {
		cli.printChain()
	}
	return nil
}

func (cli *CLI) addBlock(data string) {
	block := cli.bc.AddBlock(data)
	fmt.Printf("Block with hash %x added to blockchain", block.Hash)
}

func (cli *CLI) printChain() {
	iter := cli.bc.Iterator()

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
