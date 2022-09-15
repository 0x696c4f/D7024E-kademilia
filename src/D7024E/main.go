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
	fmt.Println(msg+"\n\n",helpText)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		printHelpExit("No command supplied.")
	}
	port := 4000
	switch os.Args[1] {
		case "start" : {
			// start the network
			if len(os.Args) >= 3 {
				var err error 
				port,err = strconv.Atoi(os.Args[2])
				if err != nil {
					printHelpExit("Invalid port.")
				}
			}
			fmt.Println("starting network on port",port)
		}
		case "join" : {
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
				port,err = strconv.Atoi(os.Args[3])
				if err != nil {
					printHelpExit("Invalid port.")
				}
			}
			fmt.Println("joining via",ip,":",port)
		}
		case "get" : {
			hash := os.Args[2]
			fmt.Println("getting ",hash)
		}
		case "put" : {
			data := os.Args[2]
			fmt.Println("storing ",data)
		}
		default:
			printHelpExit("Invalid command.")
	}

	fmt.Println("start")

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
