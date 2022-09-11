package main

import (
	"fmt"
)

func main() {
	fmt.Println("start stop ")
	TestPing()
}

func TestPing() {

	//network struct create
	net := NewNetwork()

	// kademlia struct created, one ID and one distance
	//Get the correct and accurat ID TODO
	myId := NewRandomKademliaID()

	//------------------------------
	//Get the IP correct address TODO
	//------------------------------

	//create a contact
	contactFirst := NewContact(myId, "172.0.0.2 8080") //IP address TODO

	//call ping message in network SendPingMessage(contact)
	net.SendPingMessage(&contactFirst)

}
func StartNetwork() {

	//contact := &Contact{}
	//contact.ID =

	//When starting the network connection, have the ip and the port of one node

	//get the id from the known node
	//store that Id in the bucket

}

func JoinNetwork() {

}
