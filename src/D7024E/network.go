package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type MessageBody struct {
	ContactList []Contact //shortList which is sent back
	TargetID    *KademliaID
	Data        []byte
	Hash        string
}

type Network struct {
	Node *Kademlia

	Mu          sync.Mutex
	StoreValues map[string][]byte //Store data that is recived from the store RPC
	TTLs        map[string]time.Time

	Refresh map[string]([]Contact) // map of hash -> node id to send refresh to, list of KademliaIDs
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
	UDPaddress.IP = nil
	Conn, listeningError := net.ListenUDP("udp", &UDPaddress)

	if listeningError != nil {
		return
	}
	//2------------------
	defer Conn.Close()
	//3------------------
	for {
		//4-------------------
		buffert := make([]byte, 3048)
		fmt.Println("Listening")
		step, readAddress, _ := Conn.ReadFromUDP(buffert)
		//5-------------------
		message := ByteToPacket(buffert[0:step])
		fmt.Println("recived: ", message.RPC)
		//6-------------------
		if message.RPC != "local_get" && message.RPC != "local_put" {
			fmt.Println("[NETWORK] adding contact for RPC type", message.RPC)
			network.AddContact(*message.SendingContact)
		}
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
	value := make(map[string][]byte) //store test
	network.StoreValues = value      //store test
	ttls := make(map[string]time.Time)
	network.TTLs = ttls
	refresh := make(map[string][]Contact)
	network.Refresh = refresh
	kad := NewKademlia(localIP)
	network.Node = &kad
	return network
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

	if response.RPC != "local_get" && response.RPC != "local_put" && response.RPC != "local_forget" {
		fmt.Println("[NETWORK] adding contact for RPC type", response.RPC)
		network.AddContact(*response.SendingContact)
	}

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
	} else if version == "find_Value" {
		pack.RPC = "find_Value"
		return
	} else if version == "store_Value" {
		pack.RPC = "store_Value"
		return
	}

	return Packet{}
}

func (network *Network) SendLocalGet(hash string) []byte {
	var pack = Packet{
		SendingContact: &network.Node.RoutingTable.me,
		RPC:            "local_get",
		Message: MessageBody{
			TargetID: NewKademliaID(hash),
		},
	}

	fmt.Println("[NETWORK] send local get to ", fmt.Sprintf("127.0.0.1:%d", Port))
	instance := NewContact(NewRandomKademliaID(), fmt.Sprintf("127.0.0.1:%d", Port))
	response, err := network.UDPConnectionHandler(&instance, pack)
	if err == nil {
		network.ResponseHandler(response)
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
	return response.Message.Data
}

func (network *Network) SendLocalForget(hash string) {
	var pack = Packet{
		SendingContact: &network.Node.RoutingTable.me,
		RPC:            "local_forget",
		Message: MessageBody{
			TargetID: NewKademliaID(hash),
		},
	}

	fmt.Println("[NETWORK] send local forget to ", fmt.Sprintf("127.0.0.1:%d", Port))
	instance := NewContact(NewRandomKademliaID(), fmt.Sprintf("127.0.0.1:%d", Port))
	_, err := network.UDPConnectionHandler(&instance, pack)
	if err == nil {
		fmt.Println("success")
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (network *Network) SendRefresh(target *Contact, hash string) []byte {
	var pack = Packet{
		SendingContact: &network.Node.RoutingTable.me,
		RPC:            "refresh",
		Message: MessageBody{
			TargetID: NewKademliaID(hash),
		},
	}

	fmt.Println("[NETWORK] send refresh for " + hash)
	response, err := network.UDPConnectionHandler(target, pack)
	if err == nil {
		fmt.Println("success")
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
	return response.Message.Data
}

func (network *Network) SendLocalPut(data []byte) string {
	var pack = Packet{
		SendingContact: &network.Node.RoutingTable.me,
		RPC:            "local_put",
		Message: MessageBody{
			Data: data,
		},
	}

	fmt.Println("[NETWORK] send local put to ", fmt.Sprintf("127.0.0.1:%d", Port))
	instance := NewContact(NewRandomKademliaID(), fmt.Sprintf("127.0.0.1:%d", Port))
	response, err := network.UDPConnectionHandler(&instance, pack)
	if err == nil {
		network.ResponseHandler(response)
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
	return response.Message.TargetID.String()
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
	//fmt.Println("value contact - ", contact) //delete

	pack := network.NewPacket("find_Node")
	pack.Message = MessageBody{
		TargetID: target.ID,
	}

	response, err := network.UDPConnectionHandler(contact, pack) //TODO handle the output packet

	if err == nil {
		//fmt.Println("responce", response, " error ", err) //delete
		network.ResponseHandler(response)
	} else {
		fmt.Println(err)
	}
}

func (network *Network) SendFindDataMessage(hash string, contact Contact) []byte { // TODO
	//fmt.Println("value contact store at - ", contact, "value: ", hash) //delete
	pack := network.NewPacket("find_Value")
	pack.Message = MessageBody{
		Hash: hash,
	}

	response, err := network.UDPConnectionHandler(&contact, pack) //TODO handle the output packet

	if err == nil {
		//fmt.Println("responce", response, " error ", err) //delete
		network.ResponseHandler(response)
	} else {
		fmt.Println(err)
	}
	return response.Message.Data
}

func (network *Network) SendStoreMessage(data []byte, storeAtContact *Contact) { // TODO
	//fmt.Println("value contact store at - ", storeAtContact, "value: ", data) //delete
	pack := network.NewPacket("store_Value")
	pack.Message = MessageBody{
		Data: data,
	}
	/*
		fmt.Println("Packet info ", pack.Message)
		packToByte := PacketToByte(pack)
		fmt.Println("Packet to byte info ", packToByte)
		byteToPack := ByteToPacket(packToByte)
		fmt.Println("byte to packet info ", byteToPack.Message)
	*/
	response, err := network.UDPConnectionHandler(storeAtContact, pack)
	if err == nil {
		//fmt.Println("responce", response, " error ", err) //delete
		network.ResponseHandler(response)
	} else {
		fmt.Println(err)
	}
}
func (network *Network) ForgetOld() {
	TTL := 30
	TTLunit := "s"
	defer network.Mu.Unlock()
	ttl, err := time.ParseDuration(strconv.Itoa(TTL) + TTLunit) // TTL DEFINED HERE
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ttl)
	for {
		next := time.Now().Add(ttl) // in 30s
		network.Mu.Lock()
		for k, v := range network.TTLs {
			expires := v.Add(ttl)
			if expires.Before(time.Now()) { // time added plus ttl = expires at
				fmt.Println("Timeout for ", k)
				delete(network.TTLs, k)
				delete(network.StoreValues, k)
			} else {
				if expires.Before(next) {
					next = expires
				}
			}
		}
		network.Mu.Unlock()
		fmt.Println("Next forget check in ", next.Sub(time.Now()))
		time.Sleep(next.Sub(time.Now()))

	}
}
func (network *Network) RefreshLoop() {
	TTL := 30
	TTLunit := "s"
	defer network.Mu.Unlock()
	delay, _ := time.ParseDuration(strconv.Itoa(TTL/2) + TTLunit) // TTL DEFINED HERE
	for {
		network.Mu.Lock()
		// TODO clone network.Refresh and unlock afterwards, replace network.Refresh by local copy
		//network.Mu.Unlock()
		for k, v := range network.Refresh {
			for _, n := range v {
				fmt.Println("Sending refresh for ", k)
				network.SendRefresh(&n, k) // refresh data with hash k at node n
			}
		}
		network.Mu.Unlock()
		time.Sleep(delay)

	}
}
