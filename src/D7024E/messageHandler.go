package main

func (network *Network) MessageHandler(message *Packet) Packet {
	if message.form == "ping" {
		return network.NewResponsePacket(message)
	}
	return Packet{}
}

func (network *Network) NewResponsePacket(message *Packet) (pack Packet) {
	pack = Packet{
		form:           "ping",     //TODO should it be called the same
		ID:             message.ID, //TODO Should it have it's own or message ID
		SendingContact: network.node.routingTable.me,
	}
	return
}
