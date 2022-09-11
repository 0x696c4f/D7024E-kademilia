package main

import (
	"fmt"
)

func main() {
	fmt.Println("start stop ")

	myId := NewRandomKademliaID()
	myId2 := NewRandomKademliaID()

	fmt.Println(myId)
	fmt.Println(myId2)

	contactFirst := NewContact(myId, "172.0.0.2")

	fmt.Println(contactFirst.ID)
	//network struct create
	//two kademlia struct created, one ID and one distance
	//create a contact
	//call ping message in network SendPingMessage(contact)
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
