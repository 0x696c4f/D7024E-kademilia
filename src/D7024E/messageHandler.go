package main

import "fmt"

func (network *Network) MessageHandler(message *Packet) Packet {

	if message.RPC == "ping" {
		return network.NewPingResponsePacket(message)
	} else if message.RPC == "find_Node" {
		return network.NewFindNodeResponsePacket(message)
	}
	fmt.Println("Don't wanna see")

	return Packet{}
}

func (network *Network) NewPingResponsePacket(message *Packet) (pack Packet) {
	pack = Packet{
		RPC:            "pong",
		ID:             message.ID,
		SendingContact: network.Node.RoutingTable.Me,
	}
	return
}

func (network *Network) NewFindNodeResponsePacket(packMesssage *Packet) (pack Packet) {

	closestContacts := network.Node.RoutingTable.FindClosestContacts(packMesssage.Message.TargetID, bucketSize)

	response := MessageBody{
		ContactList: closestContacts,
	}

	pack = Packet{
		RPC:            "find_Node_res",
		ID:             packMesssage.ID,
		SendingContact: network.Node.RoutingTable.Me,
		Message:        response,
	}

	return
}
