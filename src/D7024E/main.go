package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

func printHelpExit(msg string) {
	helpText := "Usage:\nstart [port]\t\t start the first node of a kademlia network\njoin <ip> [port]\t join an existing network using the node ip:port as the entrypoint\nget <hash>\t\t get the object with hash from the network\nput <data>\t\t store data into the network\n\nThe default port always is 4000\n"
	fmt.Println(msg+"\n\n", helpText)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		printHelpExit("No command supplied.")
	}

	port := "10001"

	myIP := GetOutboundIP()
	localIP := myIP.String() + ":" + port // currentNode IP
	//bnIP := GetbnIP(myIP.String())        // bootstrapNode IP
	bnIP := "172.17.0.2:10001"

	fmt.Println("Your IP is:", localIP)

	//bnID := NewKademliaID(HashData(bnIP))
	bnContact := NewContact(nil, bnIP)

	network := NewNetwork(localIP)

	switch os.Args[1] {
	case "start":
		{

			if localIP != bnIP {
				// Join network by sending LookupContact to bootstrapNode
				network.SendPingMessage(&bnContact)
			}
			network.LookupContact(&network.Node.RoutingTable.me)

		}
	case "ping":
		{
			if len(os.Args) < 3 {
				printHelpExit("No entrypoint given.")
			}
			ipStr := os.Args[2]
			ip := net.ParseIP(ipStr)
			if ip == nil {
				printHelpExit("Invalid IP")
			}
			if len(os.Args) >= 4 {
				var err error
				if err != nil {
					printHelpExit("Invalid port.")
				}
			}

			network.TestPing(ip)
		}
	case "listen":
		{
			network.Listen()
		}
	default:
		printHelpExit("Invalid command.")
	}

	//Testing call for Listen
	network.Listen()

}

func (network *Network) TestRoutingTable() {
	TestconnectIP2 := "172.17.0.5:10001"
	testContact := NewContact(NewKademliaID(HashData(TestconnectIP2)), TestconnectIP2) //IP address TODO
	//network.node.routingTable.AddContact(testContact)

	network.Node.RoutingTable.AddContact(testContact)
}

func (network *Network) PopulateRoutingTable() {

	for n := 0; n < 20; n++ {

		TestconnectIP := "172.17.0.4:10001"
		network.Node.RoutingTable.AddContact(NewContact(NewRandomKademliaID(), TestconnectIP))
	}
}

func (network *Network) TestPing(ip net.IP) {

	//create a contact
	TestconnectIP := ip.String() + ":10001"
	contactFirst := NewContact(NewKademliaID(HashData(TestconnectIP)), TestconnectIP) //IP address TODO
	fmt.Println(contactFirst)

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

func GetbnIP(ipString string) string {
	//stringList := strings.Split(ipString, ".")
	//value := stringList[1]
	bnIP := "172." + "17" + ".0.2:10001"
	return bnIP
}
func IpPortSerialize(myIP net.IP, port int) string {
	localIP := myIP.String() + ":" + strconv.Itoa(port)
	return localIP
}
