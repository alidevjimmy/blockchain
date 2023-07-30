package main

import (
	"fmt"
	"log"

	"github.com/itchyny/base58-go"
)

func IntToHex[I int32 | int64 | int](n I) []byte {
	return []byte(fmt.Sprintf("%x", n))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func Base58Encode(text []byte) []byte {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(text)
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	return encoded
}

func Base58Decode(text []byte) []byte {
	encoding := base58.BitcoinEncoding
	decoded, err := encoding.Decode(text)
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	return decoded
}
