package main

import (
	"fmt"
)

func IntToHex[I int32 | int64 | int](n I) []byte {
	return []byte(fmt.Sprintf("%x", n))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
