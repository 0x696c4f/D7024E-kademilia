package main

import (
	"fmt"
	"time"
)

func (network *Network) MessageHandler(message Packet) Packet {

	fmt.Println("[MSGHANDLER] got ", message.RPC)
	if message.RPC == "ping" {
		return network.NewPingResponsePacket(message)
	} else if message.RPC == "find_Node" {
		return network.NewFindNodeResponsePacket(message)
	} else if message.RPC == "find_Value" {
		return network.NewFindValueResponsePacket(message)
	} else if message.RPC == "store_Value" {
		return network.NewStoreResponsePacket(message)
	} else if message.RPC == "local_get" {
		return network.NewDataPacket(message)
	} else if message.RPC == "local_put" {
		return network.NewHashPacket(message)
	}
	if message.RPC == "local_forget" {
		return network.NewForgetResponsePacket(message)
	}
	if message.RPC == "refresh" {
		return network.NewRefreshResponsePacket(message)
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

func (network *Network) NewRefreshResponsePacket(packMesssage Packet) Packet {
	network.Mu.Lock()
	defer network.Mu.Unlock()
	network.TTLs[packMesssage.Message.TargetID.String()]=time.Now()

	pack := Packet{
		RPC:            "refreshed",
		SendingContact: &network.Node.RoutingTable.me,
	}

	return pack
}

func (network *Network) NewFindValueResponsePacket(packMesssage Packet) Packet {
	//Sends the data object from the map if the hash matches a stored key
	network.Mu.Lock()
	defer network.Mu.Unlock()
	fmt.Println("message ", network.StoreValues[packMesssage.Message.Hash])
	if value, found := network.StoreValues[packMesssage.Message.Hash]; found {
		fmt.Println("The value was found!! ", string(value))
		network.TTLs[packMesssage.Message.Hash]=time.Now()
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

	fmt.Println("The value was not found!! ")
	response := MessageBody{
		ContactList: network.Node.RoutingTable.FindClosestContacts(NewKademliaID(packMesssage.Message.Hash), network.Node.Alpha),
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
	network.Mu.Lock()
	defer network.Mu.Unlock()
	network.StoreValues[hashMessageData] = message.Message.Data
	network.TTLs[hashMessageData] = time.Now()
	fmt.Println("mapList ", network.StoreValues)

	pack := Packet{
		RPC:            "store_Value",
		SendingContact: &network.Node.RoutingTable.me,
	}

	return pack
}

func (network *Network) NewHashPacket(message Packet) (pack Packet) {
	network.Store(message.Message.Data)
	response := MessageBody{
		TargetID: network.Store(message.Message.Data),
	}
	pack = Packet{
		RPC:            "hash",
		SendingContact: &network.Node.RoutingTable.me,
		Message:        response,
	}
	return
}
func (network *Network) NewForgetResponsePacket(message Packet) (pack Packet) {
	network.Forget(message.Message.TargetID.String())
	pack = Packet{
		RPC:            "forgot",
		SendingContact: &network.Node.RoutingTable.me,
	}
	return
}

func (network *Network) NewDataPacket(message Packet) (pack Packet) {
	response := MessageBody{
		Data: network.LookupData(message.Message.TargetID.String()),
	}
	pack = Packet{
		RPC:            "data",
		SendingContact: &network.Node.RoutingTable.me,
		Message:        response,
	}
	return
}
