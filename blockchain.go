package main

type Blockchain struct {
	Blocks []*Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: []*Block{
			NewGenesisBlock(),
		},
	}
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}
