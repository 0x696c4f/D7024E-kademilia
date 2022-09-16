package main

import "fmt"

func (network *Network) MessageHandler(message *Packet) Packet {

	if message.RPC == "ping" {
		return network.NewResponsePacket(message)
	}
	fmt.Println("Don't wanna see")

	return Packet{}
}

func (network *Network) NewResponsePacket(message *Packet) (pack Packet) {
	pack = Packet{
		RPC:            "pong",     //TODO should it be called the same
		ID:             message.ID, //TODO Should it have it's own or message ID
		SendingContact: network.node.routingTable.me,
	}
	return
}
