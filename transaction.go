package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
)

type Transction struct {
	ID   []byte
	VIn  []TxInput
	VOut []TxOutput
}

type TxInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string // script for unlocking prev transaction output in order to spend it
}

type TxOutput struct {
	Value        int
	ScriptPubKey string
}

func NewCoinbaseTx(to, data string) *Transction {
	if data == "" {
		data = fmt.Sprintf("Reward to %s", to)
	}
	txin := TxInput{
		Txid:      []byte{},
		Vout:      -1,
		ScriptSig: data,
	}
	txout := TxOutput{
		Value:        REWARD,
		ScriptPubKey: to,
	}
	return &Transction{
		ID:   nil,
		VIn:  []TxInput{txin},
		VOut: []TxOutput{txout},
	}
}

func (txin *TxInput) CanUnlockOutput(unlockScript string) bool {
	return txin.ScriptSig == unlockScript
}

func (txout *TxOutput) CanBeUnlockedWith(unlockScript string) bool {
	return txout.ScriptPubKey == unlockScript
}

func (tx *Transction) IsCoinBase() bool {
	return len(tx.VIn) == 0
}

func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transction {
	var inputs []TxInput
	var outputs []TxOutput

	accu, validTxs := bc.FindSpendableUTXOs(from, amount)

	if accu < amount {
		log.Panic("ERROR: Not enough balance")
	}

	for txID, tx := range validTxs {
		id, err := hex.DecodeString(txID)
		if err != nil {
			log.Panic("ERRRO: unable to decode transaction id")
		}
		for _, outIdx := range tx {
			inputs = append(inputs, TxInput{
				Txid:      id,
				Vout:      outIdx,
				ScriptSig: from,
			})
		}
	}

	outTo := TxOutput{
		Value:        amount,
		ScriptPubKey: to,
	}
	outputs = append(outputs, outTo)
	if accu > amount {
		outputs = append(outputs, TxOutput{
			Value:        accu - amount,
			ScriptPubKey: from,
		})
	}

	tx := Transction{
		ID: nil,
		VIn: inputs,
		VOut: outputs,
	}

	tx.SetID()

	return &tx
}

func (tx *Transction) SetID() {
	var inOut [][]byte
	for _, in := range tx.VIn {
		inb := bytes.Join([][]byte{
			in.Txid,
			IntToHex(in.Vout),
			[]byte(in.ScriptSig),
		}, []byte{})
		inOut = append(inOut, inb)
	}
	for _, out := range tx.VOut {
		outb := bytes.Join([][]byte{
			IntToHex(out.Value),
			[]byte(out.ScriptPubKey),
		}, []byte{})
		inOut = append(inOut, outb)
	}
	data := bytes.Join(inOut, []byte{})

	hash := sha256.Sum256(data)

	tx.ID = hash[:]
}
