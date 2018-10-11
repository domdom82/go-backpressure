package client

import (
	"bytes"
	"math/rand"
)

type Client interface {
	Run()
}

func makePayload(size int) []byte {
	payload := new(bytes.Buffer)

	for i := 0; i < size; i++ {
		payload.WriteByte(byte(rand.Intn(255)))
	}
	return payload.Bytes()
}
