package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntToHex(t *testing.T) {
	var n int64 = 12983719823719832
	hex := []byte("2e209fd7f4bd98")

	h := IntToHex(n)
	assert.Equal(t, hex, h)
}
