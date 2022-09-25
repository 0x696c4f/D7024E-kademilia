package main

import "fmt"

func (network *Network) MessageHandler(message Packet) Packet {

	if message.RPC == "ping" {
		return network.NewPingResponsePacket(message)
	} else if message.RPC == "find_Node" {
		return network.NewFindNodeResponsePacket(message)
	}
	fmt.Println("Don't wanna see")

	return Packet{
		RPC: "this is the worst case ",
	}
}

func (network *Network) NewPingResponsePacket(message Packet) (pack Packet) {
	pack = Packet{
		RPC:            "ping",
		ID:             message.ID,
		SendingContact: &network.Node.RoutingTable.me,
	}
	return
}

func (network *Network) NewFindNodeResponsePacket(packMesssage Packet) Packet {

	/*closestContacts := network.Node.RoutingTable.FindClosestContacts(packMesssage.Message.TargetID, bucketSize)

	response := MessageBody{
		ContactList: closestContacts,
	}

	pack = Packet{
		RPC:            "find_Node_res",
		ID:             packMesssage.ID,
		SendingContact: network.Node.RoutingTable.me,
		Message:        response,
	}*/

	network.Node.RoutingTable.AddContact(*packMesssage.SendingContact)

	closestContacts := make([]Contact, 0)
	closestContacts = network.Node.RoutingTable.FindClosestContacts(packMesssage.Message.TargetID, bucketSize)

	response := MessageBody{
		ContactList: closestContacts,
	}

	pack := Packet{
		RPC:            "find_Node",
		ID:             packMesssage.ID,
		SendingContact: &network.Node.RoutingTable.me,
		Message:        response,
	}

	return pack
}
