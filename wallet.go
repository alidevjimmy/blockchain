package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

var (
	addressChecksumLen = 8
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallet() *Wallet {
	private, public := newKeyPeir()
	wallet := &Wallet{
		PrivateKey: *private,
		PublicKey:  public,
	}
	return wallet
}

func newKeyPeir() (*ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Println(err)
		return nil, []byte{}
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return private, pubKey
}

func (w *Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)
	versionedPayload := append([]byte{VERSION}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)

	address := Base58Encode(fullPayload)

	return address
}

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()

	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	pubRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return pubRIPEMD160
}

func checksum(payload []byte) []byte {
	fHash := sha256.Sum256(payload)
	sHash := sha256.Sum256(fHash[:])

	return sHash[:addressChecksumLen]
}

