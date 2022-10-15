package main

import (
	"fmt"

	"testing"
)

func TestNetwork(t *testing.T) {
	fmt.Println("implement network testing")
	network := NewNetwork("localhost:8000")
	contactA := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")

	network.SendRefresh(&contactA, HashData("hejsan"))
	//network.SendLocalGet(HashData("hejsan"))
}
