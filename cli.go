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
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	createBlockchainAddress := createBlockchainCmd.String("address", "", "user wallet address")
	createBlockchainName := createBlockchainCmd.String("name", "", "blockchain name")

	printchainAddress := printchainCmd.String("address", "", "user wallet address")
	printchainName := printchainCmd.String("name", "", "blockchain name")

	getBalanceAddress := getBalanceCmd.String("address", "", "user wallet address")
	getBalanceName := getBalanceCmd.String("name", "", "blockchain name")

	sendCmdFrom := sendCmd.String("from", "", "user wallet address")
	sendCmdTo := sendCmd.String("to", "", "blockchain name")
	sendCmdName := sendCmd.String("name", "", "blockchain name")
	sendCmdAmount := sendCmd.String("amount", "", "blockchain name")

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
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			return err
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			return err
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			return err
		}
	default:
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		cli.addBlock()
	}
	if printchainCmd.Parsed() {
		cli.printChain(*printchainAddress, *printchainName)
	}
	if createBlockchainCmd.Parsed() {
		cli.createBlockchain(*createBlockchainAddress, *createBlockchainName)
	}

	if getBalanceCmd.Parsed() {
		cli.getBalance(*getBalanceAddress, *getBalanceName)
	}

	if sendCmd.Parsed() {
		amount, err := strconv.Atoi(*sendCmdAmount)
		if err != nil {
			panic("invalid amount")
		}
		cli.send(*sendCmdFrom, *sendCmdTo, *sendCmdName, amount)
	}
	return nil
}

func (cli *CLI) addBlock() {
	// block := cli.bc.AddBlock(txs)
	// fmt.Printf("Block with hash %x added to blockchain", block.Hash)
	fmt.Println("addblock command is deprecated")
}

func (cli *CLI) printChain(address, blockchainName string) {
	bc := NewBlockchain(address, blockchainName)

	defer bc.db.Close()
	iter := bc.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Prev Block Hash: %x\n", block.PrevBlockHash)
		// fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Block Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.IsValid()))
		fmt.Println("Transactions: ")
		for _, tx := range block.TXs {
			fmt.Printf("TxID: %x\n", tx.ID)
			fmt.Println("Inputs: ")
			for _, in := range tx.VIn {
				fmt.Println("ScriptSig: ", in.PubKey)
				fmt.Println("TxId: ", in.Txid)
				fmt.Println("Vout: ", in.Vout)
			}
			fmt.Println("Outputs: ")
			for _, out := range tx.VOut {
				fmt.Println("ScriptPubKey: ", out.PubKeyHash)
				fmt.Println("Value: ", out.Value)
			}
		}
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) getBalance(address, blockchainName string) {
	bc := NewBlockchain(address, blockchainName)
	defer bc.db.Close()

	UTXOs := bc.FindUTXOs(address)
	balance := 0
	for _, UTXO := range UTXOs {
		balance += UTXO.Value
	}
	fmt.Printf("Your Coin Balance is: %d\n", balance)
}

func (cli *CLI) send(from, to, blockchainName string, amount int) {
	bc := NewBlockchain(from, blockchainName)
	defer bc.db.Close()

	tx := NewUTXOTransaction(from, to, amount, bc)

	bc.AddBlock([]*Transaction{tx})

	fmt.Printf("Transfer %d from %s to %s completed successfully", amount, from, to)
}

func (cli *CLI) createBlockchain(address, name string) {
	NewBlockchain(address, name)
}
