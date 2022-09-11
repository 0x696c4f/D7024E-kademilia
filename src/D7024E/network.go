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
	//have a massage which is encoded 	"message := EncodeString("Just a PING message")"
	//send the message to the dedicated address  "go network.sendUDP("PING", contact.Address, message)"

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
