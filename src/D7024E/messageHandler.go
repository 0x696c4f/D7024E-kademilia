package main

func (network *Network) MessageHandler(message Packet) Packet {

	if message.RPC == "ping" {
		return network.NewPingResponsePacket(message)
	} else if message.RPC == "find_Node" {
		return network.NewFindNodeResponsePacket(message)
	}
	else if message.RPC == "store_Value"{
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

func(network *Network) NewStoreResponsePacket(message Packet) Packet{
	pack := Packet{
		RPC:            "store_Value",
		SendingContact: &network.Node.RoutingTable.me,
	}

	return pack
}