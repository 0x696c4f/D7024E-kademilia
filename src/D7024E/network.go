package main

import (
	"fmt"
)

type Network struct {
}

func Listen(ip string, port int) {
	// TODO
}

func NewNetwork() *Network {
	net := &Network{}
	return net
}

func (network *Network) SendPingMessage(contact *Contact) {
	fmt.Println("Sending ping message")
	//have a massage which is encoded

	//Send Message
	/*
		0: Take out the contact information from contact
		1: establish connection
			1.1 If connection fail send error message
			1.2 defer closing the connection. Making sure that it closes even on errors
		2: Write through connection
		3: set a deadline for a respons
		4: If deadline expires send error message
		5: Message was responed on, add the contact to the bucket
	*/

}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
