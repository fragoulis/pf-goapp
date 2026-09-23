package util

import (
	"math/rand"
)

var randx = rand.NewSource(42)

// RandString returns a random string of length n.
func RandString(n int) string {
	const hexDigits = "0123456789ABCDEF"
	const (
		hexIdxBits = 4                 // 4 bits to represent a hexadecimal digit.
		hexIdxMask = 1<<hexIdxBits - 1 // All 1-bits, as many as hexIdxBits.
		hexIdxMax  = 63 / hexIdxBits   // # of hexadecimal digits fitting in 63 bits.
	)

	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for hexIdxMax hexadecimal digits.
	for i, cache, remain := n-1, randx.Int63(), hexIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randx.Int63(), hexIdxMax
		}
		b[i] = hexDigits[int(cache&hexIdxMask)]
		i--
		cache >>= hexIdxBits
		remain--
	}

	return string(b)
}
