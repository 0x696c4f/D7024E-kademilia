package main

import "fmt"

func (network *Network) MessageHandler(message Packet) Packet {

	fmt.Println("[MSGHANDLER] got ",message.RPC)
	if message.RPC == "ping" {
		return network.NewPingResponsePacket(message)
	} else if message.RPC == "find_Node" {
		return network.NewFindNodeResponsePacket(message)
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
	network.Node.Store(message.Message.Data)
	response := MessageBody {
		TargetID: network.Node.Store(message.Message.Data),
	}
	pack = Packet{
		RPC:            "hash",
		SendingContact: &network.Node.RoutingTable.me,
		Message: response,
	}
	return
}

func (network *Network) NewDataPacket(message Packet) (pack Packet) {
	response := MessageBody {
		Data: network.Node.LookupData(message.Message.TargetID.String()),
	}
	pack = Packet{
		RPC:            "data",
		SendingContact: &network.Node.RoutingTable.me,
		Message: response,
	}
	return
}
