package main

import (
	"encoding/json"
	"math/rand"
	"net"
	"strconv"
)

type Packet struct {
	RPC            string
	ID             *KademliaID
	SendingContact *Contact
	Message        []byte
}

func NewPingPacket() Packet {
	packet := Packet{
		RPC: "ping",
		ID:  NewRandomKademliaID(),
	}
	return packet
}

func main() {

	//create a contact
	TestconnectIP := "172.17.0.4:8080"
	contactFirst := NewContact(NewRandomKademliaID(), TestconnectIP) //IP address TODO

	//call ping message in network SendPingMessage(contact)
	UDPConnectionHandler(&contactFirst, "hejsan")
}

func UDPConnectionHandler(contact *Contact, msgPacket string) {

	packet := Packet{
		RPC: "ping",
		ID:  NewRandomKademliaID(),
	}
	UDPaddress := GetUDPAddress(contact)
	msgMarshal := PacketToByte(packet)
	//1-------------
	Conn, _ := net.DialUDP("udp", nil, &UDPaddress)

	defer Conn.Close()
	//2-------------
	Conn.Write([]byte(msgMarshal)) //TODO check if write is correct, could be WriteToUDP

}

func PacketToByte(pkt Packet) []byte {
	message, _ := json.Marshal(pkt)
	return message
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

const IDLength = 20

func NewRandomKademliaID() *KademliaID {
	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = uint8(rand.Intn(256))
	}
	return &newKademliaID
}

type Kademlia struct {
	Me Contact //my own information
	//routingTable *RoutingTable //Everyone else infromation
}

type Contact struct {
	ID       *KademliaID
	Address  string
	distance *KademliaID
}
type KademliaID [IDLength]byte

// NewContact returns a new instance of a Contact
func NewContact(id *KademliaID, address string) Contact {
	return Contact{id, address, nil}
}
