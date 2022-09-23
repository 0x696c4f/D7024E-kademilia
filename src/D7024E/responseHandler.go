package main

import "fmt"

func (network *Network) ResponseHandler(response *Packet) {
	if response.RPC == "pong" {
		network.AddToRoutingTable(response.SendingContact)
	} else if response.RPC == "find_Node_res" {
		network.HandleFindNodeResponse(response)
	}
}

func (network *Network) HandleFindNodeResponse(response *Packet) {
	//TODO
	fmt.Println("handle the find node responce")
}
