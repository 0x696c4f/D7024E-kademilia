package main

func (network *Network) ResponseHandler(response Packet) {
	if response.RPC == "find_Node" {
		network.HandleFindNodeResponse(response)
	}
}

func (network *Network) HandleFindNodeResponse(response Packet) {
	//TODO
	network.Node.Shortlist.Append(response.Message.ContactList)

}
