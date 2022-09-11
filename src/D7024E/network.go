package main

import (
	"fmt"
)

type Network struct {
	//TODO
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
	//have a massage which is encoded or set up in some way

	//Send Message
	/*
		0: Take out the contact information from contact
			0.1: Get UDP adress
		1: establish connection (https://pkg.go.dev/net)
			Answer: DialUDP(network string, laddr, raddr *UDPAddr)
			1.1 If connection fail send error message
			1.2 defer closing the connection. Making sure that it closes even on errors
		2: Write through connection
		3: set a deadline for a respons
			Answer: SetDeadline(t time.Time) error
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
