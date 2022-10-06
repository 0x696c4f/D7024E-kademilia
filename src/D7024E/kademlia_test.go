package main

import (
	"fmt"
	"testing"
)

func TestKademlia(t *testing.T) {
	fmt.Println("implement kademlia testing")
	network := NewNetwork("localhost:8000")
	network.LookupContact(&network.Node.RoutingTable.me)

	contactA := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	contactB := NewContact(NewKademliaID("EEEEEEEE00000000000000000000000000000000"), "localhost:8001")
	contactC := NewContact(NewKademliaID("DDDDDDDD00000000000000000000000000000000"), "localhost:8001")
	contactD := NewContact(NewKademliaID("CCCCCCCC00000000000000000000000000000000"), "localhost:8001")
	contactE := NewContact(NewKademliaID("BBBBBBBB00000000000000000000000000000000"), "localhost:8001")
	contactF := NewContact(NewKademliaID("AAAAAAAA00000000000000000000000000000000"), "localhost:8001")

	network.AddContact(contactA)
	network.LookupContact(&network.Node.RoutingTable.me)

	network.AddContact(contactB)
	network.AddContact(contactC)
	network.AddContact(contactD)
	network.AddContact(contactE)
	network.AddContact(contactF)
	network.LookupContact(&network.Node.RoutingTable.me)
	network.LookupData(HashData("hejsan"))
	//	network.LookupData("hejsan") //will fail
	b := []byte("ABC€")

	network.Store(b)
	fmt.Println("-------------------------")
	fmt.Println("")
}
