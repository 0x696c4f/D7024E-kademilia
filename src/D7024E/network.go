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
	Node *Kademlia

	storeValues map[KademliaID]string //Store data that is recived from the store RPC
}

type Packet struct {
	RPC            string
	SendingContact *Contact
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
func (network *Network) Listen() {
	// TODO
	//1------------------
	UDPaddress := GetUDPAddress(&network.Node.RoutingTable.me)
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
		network.AddContact(*message.SendingContact)
		//7-------------------
		response := network.MessageHandler(message)
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
	kad := NewKademlia(localIP)
	network.Node = &kad
	storeValues := make(map[KademliaID]string)
	return network, storeValues
}

func (network *Network) JoinNetwork(contactKnown *Contact) {
	network.AddContact(*contactKnown)
	network.LookupContact(&network.Node.RoutingTable.me)
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
	Conn.Write([]byte(msgMarshal))
	//3-------------
	Conn.SetDeadline(time.Now().Add(100 * time.Millisecond))
	//4-------------
	buffert := make([]byte, 20000)
	step, _, readError := Conn.ReadFromUDP(buffert)
	//5-------------
	response := ByteToPacket(buffert[0:step])
	//6-------------
	if readError != nil {
		//TODO  delete the contact which didn't repsonde
		return Packet{}, readError
	}

	network.AddContact(*response.SendingContact)

	return response, nil

}

func (network *Network) NewPacket(version string) (pack Packet) {
	pack = Packet{
		SendingContact: &network.Node.RoutingTable.me,
	}
	if version == "ping" {
		pack.RPC = "ping"
		return
	} else if version == "find_Node" {
		pack.RPC = "find_Node"
		return
	}else if version == "store_Value"{
		pack.RPC = "store_Value"
		return
	}

	return Packet{}
}

func (network *Network) SendPingMessage(contact *Contact) (Packet, error) {

	response, err := network.UDPConnectionHandler(contact, network.NewPacket("ping"))
	if err == nil {
		network.ResponseHandler(response)
	} else {
		fmt.Println(err)
	}

	return response, err
}

func (network *Network) SendFindContactMessage(contact *Contact, target *Contact) {
	fmt.Println("value contact - ", contact) //delete

	pack := network.NewPacket("find_Node")
	pack.Message = MessageBody{
		TargetID: target.ID,
	}

	response, err := network.UDPConnectionHandler(contact, pack) //TODO handle the output packet

	if err == nil {
		fmt.Println("responce", response, " error ", err) //delete
		network.ResponseHandler(response)
	} else {
		fmt.Println(err)
	}
}

func (network *Network) SendFindDataMessage(hash string) {// TODO
	
}

func (network *Network) SendStoreMessage(data []byte, storeAtContact *Contact) {// TODO
	pack := network.NewPacket("store_Value")
	pack.Message = MessageBody{
		TargetID: storeAtContact.ID,
	}
	
	response, err := network.UDPConnectionHandler(storeAtContact, pack)
	if err == nil {
		network.ResponseHandler(response)
	} else {
		fmt.Println(err)
	}

	return response, err
}
