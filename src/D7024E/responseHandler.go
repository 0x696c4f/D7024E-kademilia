package main

func (network *Network) ResponseHandler(response *Packet) {
	if response.RPC == "ping" {
		network.AddToRoudingTable(&response.SendingContact)
	}
}
