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

	if localIP != "172.17.0.2:8080" { //TODO set up a universal first IP address ending with 0.2:8080
		JoinNetwork()
	}
	//net.TestPing(localIP)

	/*
		-:Check if this container is the first one in the system
			-: If it's not the first then join the existing network through messageing the first node
				-: Has IP Address ending with 0.2
		-:Start a listiner
	*/

}

func JoinNetwork() {
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
