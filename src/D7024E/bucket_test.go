package main

import (
	"fmt"
	"testing"
)

func TestBucket(t *testing.T) {
	fmt.Println("implement bucket testing")
	network := NewNetwork("localhost:8000")

	for i := 0; i < 50; i++ {
		network.AddContact(NewContact(NewRandomKademliaID(), "localhost:8001"))
	}

	fmt.Println("-------------------------")
	fmt.Println("")
}
