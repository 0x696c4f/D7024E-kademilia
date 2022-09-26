package main

import "fmt"

func (network *Network) ResponseHandler(response Packet) {
	fmt.Println("You got a ", response.RPC, " response")
	if response.RPC == "find_Node" {
		network.HandleFindNodeResponse(response)
	}
}

func (network *Network) HandleFindNodeResponse(response Packet) {
	//TODO
	network.Node.Shortlist.Append(response.Message.ContactList)
}
