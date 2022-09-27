package main

func (network *Network) MessageHandler(message Packet) Packet {

	if message.RPC == "ping" {
		return network.NewPingResponsePacket(message)
	} else if message.RPC == "find_Node" {
		return network.NewFindNodeResponsePacket(message)
	}
	if message.RPC == "local_get" {
		network.Node.LookupData(message.Message.TargetID.String())
	}
	if message.RPC == "local_put" {
		network.Node.Store(message.Message.Data)
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
