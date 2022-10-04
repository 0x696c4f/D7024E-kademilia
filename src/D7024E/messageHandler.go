package main

import "fmt"

func (network *Network) MessageHandler(message Packet) Packet {

	if message.RPC == "ping" {
		return network.NewPingResponsePacket(message)
	} else if message.RPC == "find_Node" {
		return network.NewFindNodeResponsePacket(message)
	} else if message.RPC == "find_Value" {
		return network.NewFindValueResponsePacket(message)
	} else if message.RPC == "store_Value" {
		return network.NewStoreResponsePacket(message)
	}

	return Packet{}
}

func (network *Network) NewPingResponsePacket(message Packet) (pack Packet) {
	pack = Packet{
		RPC:            "ping",
		SendingContact: &network.Node.RoutingTable.me,
	}
	return
}

func (network *Network) NewFindNodeResponsePacket(packMesssage Packet) Packet {

	response := MessageBody{
		ContactList: network.Node.RoutingTable.FindClosestContacts(packMesssage.Message.TargetID, network.Node.Alpha),
	}

	pack := Packet{
		RPC:            "find_Node",
		SendingContact: &network.Node.RoutingTable.me,
		Message:        response,
	}

	return pack
}

func (network *Network) NewFindValueResponsePacket(packMesssage Packet) Packet {
	//Sends the data object from the map if the hash matches a stored key
	var value []byte
	value = network.StoreValues[packMesssage.Message.Hash]

	if value != nil {
		fmt.Println("The value was found!! ", value)
		response := MessageBody{
			Data: value,
		}

		pack := Packet{
			RPC:            "find_Value",
			SendingContact: &network.Node.RoutingTable.me,
			Message:        response,
		}
		return pack
	}

	fmt.Println("The value was not found!! ", value)
	response := MessageBody{
		ContactList: network.Node.RoutingTable.FindClosestContacts(packMesssage.Message.TargetID, network.Node.Alpha),
	}

	pack := Packet{
		RPC:            "find_Value",
		SendingContact: &network.Node.RoutingTable.me,
		Message:        response,
	}

	return pack
}

func (network *Network) NewStoreResponsePacket(message Packet) Packet {
	fmt.Println("what is the message data: ", message.Message.Data)
	hashMessageData := HashData(string(message.Message.Data))
	//valueID := NewKademliaID(hashMessageData)
	fmt.Println("data to be stored: ", hashMessageData)
	//fmt.Println("what is the new kademliaID: ", valueID)
	network.StoreValues[hashMessageData] = message.Message.Data
	fmt.Println("mapList ", network.StoreValues)

	pack := Packet{
		RPC:            "store_Value",
		SendingContact: &network.Node.RoutingTable.me,
	}

	return pack
}
