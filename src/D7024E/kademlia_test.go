package main

import (
	"fmt"
	"testing"
)

func TestKademlia(t *testing.T) {
	fmt.Println("implement kademlia testing")
	network := NewNetwork("localhost:8000")

	network.LookupContact(&network.Node.RoutingTable.me)

	fmt.Println("-------------------------")
	fmt.Println("")
}
