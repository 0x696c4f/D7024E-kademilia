package main

func (network *Network) ResponseHandler(response *Packet) {
	if response.RPC == "pong" {
		network.AddToRoutingTable(response.SendingContact)
	} else if response.RPC == "find_Node_res" {
		network.HandleFindNodeResponse(response)
	}
}

func (network *Network) HandleFindNodeResponse(response *Packet) {
	//TODO
	network.AddToRoutingTable(response.SendingContact)

	network.Node.Shortlist.Lock()
	network.Node.Shortlist.Append(response.Message.ContactList)
	network.Node.Shortlist.Unlock()

}
