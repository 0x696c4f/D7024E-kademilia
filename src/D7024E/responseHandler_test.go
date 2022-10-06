package main

import (
	"fmt"
	"testing"
)

func TestResponseHandler(t *testing.T) {
	fmt.Println("implement responseHandler testing")
	fmt.Println("-------------------------")
	network := NewNetwork("localhost:8000")

	messageEmpty := network.NewPacket("empty")

	network.ResponseHandler(messageEmpty)
	fmt.Println("")
}
