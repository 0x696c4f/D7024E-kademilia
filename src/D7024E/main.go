package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Println("start")
	port := "8080"

	myIP := GetOutboundIP()
	localIP := myIP.String() + ":" + port

	network := NewNetwork()
	network.Node = NewKademlia(localIP)

	gatewayIP := GetGatewayIP()

	if localIP != gatewayIP {
		knownContact := NewContact(NewKademliaID(HashData(gatewayIP)), gatewayIP)
		JoinNetwork(&knownContact)
	}
	network.TestPing()

	//correct way to call listening
	//go network.Listen() //why we use go https://www.golang-book.com/books/intro/10

	//Testing call for Listen
	network.Listen()
}

func (network *Network) TestPing() {

	//create a contact
	TestconnectIP := "172.17.0.4:8080"
	contactFirst := NewContact(NewKademliaID(HashData(TestconnectIP)), TestconnectIP) //IP address TODO

	//call ping message in network SendPingMessage(contact)
	network.SendPingMessage(&contactFirst)
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func GetGatewayIP() (gatewayIP string) { //TODO set up a universal first IP address ending with xxx.xxx.xxx.2:8080
	gatewayIP = "172.17.0.2:8080"
	return
}
