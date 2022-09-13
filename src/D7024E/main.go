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

	net := NewNetwork()
	net.Node = NewKademlia(localIP)

	startNodeIP := GetStartNodeIP()

	if localIP != startNodeIP {
		knownContact := NewContact(NewKademliaID(HashData(startNodeIP)), startNodeIP)
		JoinNetwork(&knownContact)
	}
	//net.TestPing(localIP)

	Listen()

}

func GetStartNodeIP() (startNode string) { //TODO set up a universal first IP address ending with 0.2:8080
	startNode = "172.17.0.2:8080"
	return
}

func JoinNetwork(contactKnown *Contact) {
	fmt.Println("join network")
	//TODO
}

func (net *Network) TestPing(ip string) {

	//create a contact
	TestconnectIP := "172.17.0.3:8080"
	contactFirst := NewContact(NewKademliaID(HashData(TestconnectIP)), TestconnectIP) //IP address TODO

	//call ping message in network SendPingMessage(contact)
	net.SendPingMessage(&contactFirst)
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
