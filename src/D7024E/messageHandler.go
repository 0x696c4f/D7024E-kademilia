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
	}
	if message.RPC == "local_get" {
		return network.NewDataPacket(message)
	}
	if message.RPC == "local_put" {
		return network.NewHashPacket(message)
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
	network.Mu.Lock()
	defer network.Mu.Unlock()
	fmt.Println("message ", network.StoreValues[packMesssage.Message.Hash])
	if value, found := network.StoreValues[packMesssage.Message.Hash]; found {
		fmt.Println("The value was found!! ", string(value))
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
	ttl,_ :=time.ParseDuration("30s") // TTL
	network.TTLs[hashMessageData] = time.Now().Add(ttl)
	fmt.Println("mapList ", network.StoreValues)

	pack := Packet{
		RPC:            "store_Value",
		SendingContact: &network.Node.RoutingTable.me,
	}

	return pack
}

func (network *Network) NewLocalGetPacket(message Packet) (pack Packet) {
	pack = Packet{
		RPC:            "local_get",
		SendingContact: &network.Node.RoutingTable.me,
	}
	return
}

func (network *Network) NewLocalPutPacket(message Packet) (pack Packet) {
	pack = Packet{
		RPC:            "local_put",
		SendingContact: &network.Node.RoutingTable.me,
	}
	return
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
