package main

func (network *Network) ResponseHandler(response *Packet) {
	if response == "ping" {
		AddToRoudingTable(response.SendingContact)
	}
}
