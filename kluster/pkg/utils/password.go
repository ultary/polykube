package utils

import (
	"bytes"
	"math/rand/v2"
)

func NewPassword() string {

	const (
		letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		specialBytes = "!@#$%^&*()_+-=[]{}\\|;':\",.<>/?`~"
		numBytes     = "0123456789"
	)

	var buf bytes.Buffer
	for _ = range 3 {
		c := letterBytes[rand.IntN(len(letterBytes))]
		buf.WriteByte(c)
	}
	for _ = range 3 {
		c := specialBytes[rand.IntN(len(specialBytes))]
		buf.WriteByte(c)
	}
	for _ = range 3 {
		c := numBytes[rand.IntN(len(numBytes))]
		buf.WriteByte(c)
	}
	for _ = range 7 {
		mixedBytes := letterBytes + specialBytes + numBytes
		c := mixedBytes[rand.IntN(len(mixedBytes))]
		buf.WriteByte(c)
	}
	shuffled := buf.Bytes()
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return string(shuffled)
}
