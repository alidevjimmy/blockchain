package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

type Transaction struct {
	ID   []byte
	VIn  []TxInput
	VOut []TxOutput
}

type TxInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}

type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to %s", to)
	}
	txin := TxInput{
		Txid:   []byte{},
		Vout:   -1,
		PubKey: []byte{},
	}
	txout := TxOutput{
		Value:      REWARD,
		PubKeyHash: []byte{},
	}
	return &Transaction{
		ID:   nil,
		VIn:  []TxInput{txin},
		VOut: []TxOutput{txout},
	}
}

func (txin *TxInput) CanUnlockOutput(unlockScript []byte) bool {
	return bytes.Compare(txin.PubKey, unlockScript) == 0
}

func (txout *TxOutput) CanBeUnlockedWith(unlockScript []byte) bool {
	return bytes.Compare(txout.PubKeyHash, unlockScript) == 0
}

func (tx *Transaction) IsCoinBase() bool {
	return len(tx.VIn) == 0
}

func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
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
				Txid:   id,
				Vout:   outIdx,
				PubKey: []byte{},
			})
		}
	}

	outTo := TxOutput{
		Value:      amount,
		PubKeyHash: []byte{},
	}
	outputs = append(outputs, outTo)
	if accu > amount {
		outputs = append(outputs, TxOutput{
			Value:      accu - amount,
			PubKeyHash: []byte{},
		})
	}

	tx := Transaction{
		ID:   nil,
		VIn:  inputs,
		VOut: outputs,
	}

	tx.SetID()

	return &tx
}

func (tx *Transaction) SetID() {
	var inOut [][]byte
	for _, in := range tx.VIn {
		inb := bytes.Join([][]byte{
			in.Txid,
			IntToHex(in.Vout),
			[]byte(in.PubKey),
		}, []byte{})
		inOut = append(inOut, inb)
	}
	for _, out := range tx.VOut {
		outb := bytes.Join([][]byte{
			IntToHex(out.Value),
			[]byte(out.PubKeyHash),
		}, []byte{})
		inOut = append(inOut, outb)
	}
	data := bytes.Join(inOut, []byte{})

	hash := sha256.Sum256(data)

	tx.ID = hash[:]
}

func (txin *TxInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(txin.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

func (txout *TxOutput) Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-8]
	txout.PubKeyHash = pubKeyHash
}

func (txout *TxOutput) IsLockedWith(pubKeyHash []byte) bool {
	return bytes.Compare(pubKeyHash, txout.PubKeyHash) == 0
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinBase() {
		return
	}

	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.VIn {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.VIn[inID].Signature = nil
		txCopy.VIn[inID].PubKey = prevTx.VOut[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.VIn[inID].PubKey = nil

		r, s, _ := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		signature := append(r.Bytes(), s.Bytes()...)

		tx.VIn[inID].Signature = signature
	}
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, vin := range tx.VIn {
		inputs = append(inputs, TxInput{vin.Txid, vin.Vout, nil, nil})
	}

	for _, vout := range tx.VOut {
		outputs = append(outputs, TxOutput{vout.Value, vout.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.VIn {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.VIn[inID].Signature = nil
		txCopy.VIn[inID].PubKey = prevTx.VOut[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.VIn[inID].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}

	return true
}
