package main

func (network *Network) ResponseHandler(response *Packet) {
	if response.RPC == "pong" {
		network.AddToRoutingTable(response.SendingContact)
	}
}
