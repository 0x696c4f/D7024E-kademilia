package main

import (
	"fmt"
	"net"
	"time"
)

type MessageBody struct {
	ContactList []Contact //shortList which is sent back
	TargetID    *KademliaID
}

type Network struct {
	Node Kademlia
}

type Packet struct {
	RPC            string
	ID             *KademliaID //TODO delete
	SendingContact Contact
	Message        MessageBody
}

/*
	1:Open UDPPort for it to listen in on.
		1.1:What UDP address should we listen in on
	2:Close the connection
		answer: defer connection.Close()
	3:Create for loop to handle the inputs
	4:Read the input
	5:convert into unmarshaldata
	6:add contact to the bucket
	7:handle inquary
	8:Unmarhal data
	9:send back respons
*/
func (network *Network) Listen( /*ip string, port int*/ ) {
	// TODO
	//1------------------
	UDPaddress := GetUDPAddress(&network.Node.RoutingTable.Me)
	Conn, listeningError := net.ListenUDP("udp", &UDPaddress)

	if listeningError != nil {
		return
	}
	//2------------------
	defer Conn.Close()
	//3------------------
	for {
		//4-------------------
		buffert := make([]byte, 1000)
		fmt.Println("Listening")
		step, readAddress, _ := Conn.ReadFromUDP(buffert)
		//5-------------------
		message := ByteToPacket(buffert[0:step])
		fmt.Println("recived: ", message.RPC)
		//6-------------------
		network.AddToRoutingTable(message.SendingContact)
		//7-------------------
		response := network.MessageHandler(&message)
		//8-------------------
		responsMarshal := PacketToByte(response)
		//9-------------------
		_, respondError := Conn.WriteToUDP([]byte(responsMarshal), readAddress)
		if respondError != nil {
			fmt.Println(respondError)
		}
	}
}

func NewNetwork(localIP string) *Network {
	network := &Network{}
	network.Node = network.NewKademlia(localIP)
	return network
}

func (network *Network) JoinNetwork(contactKnown *Contact) {
	network.AddToRoutingTable(*contactKnown)
	network.Node.LookupContact(&network.Node.RoutingTable.Me) //TODO Most move back to the past alt network.node.routingTable.me
}

func (network *Network) UDPConnectionHandler(contact *Contact, msgPacket Packet) (Packet, error) {
	UDPaddress := GetUDPAddress(contact)

	msgMarshal := PacketToByte(msgPacket)

	//1-------------
	Conn, dialError := net.DialUDP("udp", nil, &UDPaddress)
	if dialError != nil {
		return Packet{}, dialError
	}

	defer Conn.Close()
	//2-------------
	Conn.Write([]byte(msgMarshal)) //TODO check if write is correct, could be WriteToUDP
	//3-------------
	Conn.SetDeadline(time.Now().Add(100 * time.Millisecond))
	//4-------------
	buffert := make([]byte, 1000)
	step, _, readError := Conn.ReadFromUDP(buffert)
	//5-------------
	response := ByteToPacket(buffert[0:step])
	//6-------------
	if readError != nil {
		//TODO  delete the contact which didn't repsonde
		return Packet{}, readError
	}
	return response, nil

}

func (network *Network) NewPacket(version string) (pack Packet) {
	if version == "ping" {
		pack = Packet{
			RPC:            "ping",
			ID:             NewRandomKademliaID(),
			SendingContact: network.Node.RoutingTable.Me,
		}
	} else if version == "find_Node" {
		pack = Packet{
			RPC:            "find_Node",
			ID:             NewRandomKademliaID(),
			SendingContact: network.Node.RoutingTable.Me,
		}
	}
	return
}

func (network *Network) SendPingMessage(contact *Contact) (Packet, error) {

	response, err := network.UDPConnectionHandler(contact, network.NewPacket("ping")) //TODO handle the output packet
	if err != nil {
		fmt.Println(err)
		return response, err
	} else {
		network.ResponseHandler(&response)
	}

	return response, nil
}

func (network *Network) SendFindContactMessage(contact *Contact, target *Contact) {
	fmt.Println("value me", network.Node.RoutingTable.Me)
	message := MessageBody{
		TargetID: target.ID,
	}
	pack := network.NewPacket("find_Node")
	pack.Message = message
	response, err := network.UDPConnectionHandler(contact, pack) //TODO handle the output packet

	if err != nil {
		fmt.Println(err)
		//return response, err
	} else {
		network.ResponseHandler(&response)
	}
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
