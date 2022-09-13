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
	TestPing(localIP)
}

func TestPing(ip string) {

	//network struct create
	net := NewNetwork()
	net.Node = NewKademlia(ip)

	//create a contact
	TestconnectIP := "172.17.0.7:8080"

	contactFirst := NewContact(NewKademliaID(HashData(TestconnectIP)), TestconnectIP) //IP address TODO

	fmt.Println(contactFirst.Address)
	//call ping message in network SendPingMessage(contact)
	//net.SendPingMessage(&contactFirst)

}

// (https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go)
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func StartNetwork() {

	//contact := &Contact{}
	//contact.ID =

	//When starting the network connection, have the ip and the port of one node

	//get the id from the known node
	//store that Id in the bucket

}

func JoinNetwork() {

}
