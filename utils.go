package main

import "fmt"

func IntToHex[I int32 | int64](n I) []byte {
	return []byte(fmt.Sprintf("%x", n))
}