package main

import (
	"encoding/json"
	"time"

	"fmt"
	"net"
	"strconv"
)

type Network struct {
	Node Kademlia
}

type Packet struct {
	RPC            string
	ID             *KademliaID
	SendingContact *Contact
	Message        []byte
}

func Listen( /*ip string, port int*/ ) {
	// TODO
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
}

func NewNetwork() *Network {
	net := &Network{}
	return net
}

func (network *Network) NewPingPacket() Packet {
	packet := Packet{
		RPC:            "ping",
		ID:             NewRandomKademliaID(),
		SendingContact: &network.Node.Me,
	}
	return packet
}

func (network *Network) SendPingMessage(contact *Contact) {
	fmt.Println("Sending ping message")

	_, err := network.UDPConnectionHandler(contact, network.NewPingPacket()) //TODO handle the output packet
	if err != nil {
		fmt.Println("------------Error-------------")
		fmt.Println(err)
		fmt.Println("------------------------------")
	} else {
		fmt.Println("Success")
	}
}

//Send Message
/*
	0: Take out the contact information from contact
		0.1: Get UDP adress by etiblating an UDPAddr
		Answer: UDPaddress := GetUDPAddress(contact)
	1: establish connection (https://pkg.go.dev/net)
		Answer: Conn, err := DialUDP(network string, laddr, raddr *UDPAddr)
				if err != nil{
					//return -- some problem has occured}
		1.1 If connection fail send error message (https://go.dev/doc/tutorial/handle-errors)
		1.2 defer closing the connection. Making sure that it closes even on errors
	2: Write through connection
		Answer: Write(b []byte)   alt:   int , err := Write(b []byte)
	3: set a deadline for a respons (https://github.com/golang/go/issues/14490)
		Answer: Conn.SetDeadline(time.Now().Add(time.Secound)) error
	4: Check for the response through the UDP connection
		Answer:	 ReadFromUDP(b []byte) (n int, addr *UDPAddr, err error)
	5: Unmarshel the incoming message
	6: If deadline expires send error message
	7: Message was responed on, add the contact to the bucket
*/
func (network *Network) UDPConnectionHandler(contact *Contact, msgPacket Packet) (Packet, error) {

	UDPaddress := GetUDPAddress(contact)
	msgMarshal := PacketToByte(msgPacket) //convert Packet to []byte

	//1-------------
	Conn, dialError := net.DialUDP("udp", nil, &UDPaddress)
	if dialError != nil {
		return Packet{}, dialError
	}
	defer Conn.Close()

	//2-------------
	Conn.Write(msgMarshal)

	//3-------------
	Conn.SetDeadline(time.Now().Add(100 * time.Millisecond)) //TODO set the time 1 to something else resonable

	//4-------------
	buffert := make([]byte, 1000)
	step, _, readError := Conn.ReadFromUDP(buffert) // response, remoteAddr, readError //TODO handle the responce

	//5-------------
	response := ByteToPacket(buffert[0:step])

	//6-------------
	if readError != nil {
		return Packet{}, readError
	}

	return response, nil

}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func PacketToByte(pkt Packet) []byte {
	message, _ := json.Marshal(pkt)
	return message
}
func ByteToPacket(message []byte) Packet {
	pkt := Packet{}
	json.Unmarshal(message, &pkt)
	return pkt
}

func GetUDPAddress(contact *Contact) net.UDPAddr {
	addr, port, _ := net.SplitHostPort(contact.Address) // "_" is a possible error message
	IPAddr := net.ParseIP(addr)
	intPort, _ := strconv.Atoi(port) //https://www.golangprograms.com/how-to-convert-string-to-integer-type-in-go.html

	UDPaddress := net.UDPAddr{
		IP:   IPAddr,
		Port: intPort,
		//Zone string // IPv6 scoped addressing zone
	}
	return UDPaddress
}
