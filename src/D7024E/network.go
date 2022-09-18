package main

import (
	"fmt"
	"net"
	"time"
)

type Network struct {
	node Kademlia

	//Trying somethings
	candidateListCh chan []Contact    //channel where candidateLists from the go routines will be written to
	inactiveNodes   ContactCandidates //used in kademlia.go
}

type Packet struct {
	RPC            string
	ID             *KademliaID
	SendingContact Contact
	Message        []byte
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
	UDPaddress := GetUDPAddress(&network.node.routingTable.me)
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

func (network *Network) NewPacket(version string) (pack Packet) {
	if version == "ping" {
		pack = Packet{
			RPC:            "ping",
			ID:             NewRandomKademliaID(),
			SendingContact: network.node.routingTable.me,
		}
	}
	return
}

func NewNetwork(localIP string) *Network {
	network := &Network{}
	network.node = NewKademlia(localIP)
	return network
}

func JoinNetwork(contactKnown *Contact) {
	fmt.Println("join network")
	//TODO
}

func (network *Network) SendPingMessage(contact *Contact) (Packet, error) {
	response, err := network.UDPConnectionHandler(contact, network.NewPacket("ping")) //TODO handle the output packet

	if err != nil {
		fmt.Println(err)
		return response, err
	}

	fmt.Println(response.RPC)
	network.ResponseHandler(&response)

	return response, nil
}

func (network *Network) UDPConnectionHandler(contact *Contact, msgPacket Packet) (Packet, error) {
	UDPaddress := GetUDPAddress(contact)
	msgMarshal := PacketToByte(msgPacket)
	//1-------------
	Conn, dialError := net.DialUDP("udp", nil, &UDPaddress)
	if dialError != nil {
		fmt.Println("test")
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
		return Packet{}, readError
	}
	return response, nil

}

// Sends a FIND_NODE_RPC to a contact
func (network *Network) SendFindContactMessage(contact *Contact, target *Contact) {
	// TODO
	UDPaddress := GetUDPAddress(contact)
	Conn, dialError := net.DialUDP("udp", nil, &UDPaddress)
	if dialError != nil {
		fmt.Println("test")

	}
	targetString := target.String()
	//Send FIND_NODE rpc together with the target contact
	fmt.Fprintf(Conn, "FIND_NODE_RPC;"+targetString+";"+network.node.routingTable.me.ID.String()+"\n")
	Conn.Close()
}

// Sends a FIND_VALUE_RPC to a contact
func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

// Sends a STORE_VALUE_RPC to a contact
func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
