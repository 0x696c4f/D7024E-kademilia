package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

var Port int

func printHelpExit(msg string) {
	helpText := "Usage:\nstart [port]\t\t start the first node of a kademlia network\njoin <ip> [port]\t join an existing network using the node ip:port as the entrypoint\nget <hash>\t\t get the object with hash from the network\nput <data>\t\t store data into the network\n\nThe default port always is 4000\n"
	fmt.Println(msg+"\n\n", helpText)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		printHelpExit("No command supplied.")
	}

	Port = 8080
	myIP := GetOutboundIP()
	localIP := IpPortSerialize(myIP, Port)

	network := NewNetwork(localIP)

	normal := true
	switch os.Args[1] {
	case "start":
		{
			// start the network
			if len(os.Args) >= 3 {
				var err error
				Port, err = strconv.Atoi(os.Args[2])
				if err != nil {
					printHelpExit("Invalid port.")
				}
			}
			fmt.Println("starting network on port", Port)
			if len(os.Args) >= 4 {
				if os.Args[3] == "test" {
					normal = false
				}
			}
		}
	case "join":
		{
			remoteport := Port
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
				remoteport, err = strconv.Atoi(os.Args[3])
				if err != nil {
					printHelpExit("Invalid port.")
				}
			}
			if len(os.Args) >= 5 {
				if os.Args[4] == "test" {
					normal = false
				}
			}

			gatewayIP := IpPortSerialize(ip, remoteport)
			fmt.Println("joining via", ip, ":", remoteport)
			knownContact := NewContact(NewKademliaID(HashData(gatewayIP)), gatewayIP)
			network.JoinNetwork(&knownContact)
		}
	case "ping":
		{
			remoteport := Port
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
				remoteport, err = strconv.Atoi(os.Args[3])
				if err != nil {
					printHelpExit("Invalid port.")
				}
			}
			if len(os.Args) >= 5 {
				if os.Args[4] == "test" {
					normal = false
				}
			}

			connectIP := IpPortSerialize(ip, remoteport)
			pingContact := NewContact(NewKademliaID(HashData(connectIP)), connectIP) //IP address TODO

			network.SendPingMessage(&pingContact)
		}
	case "get":
		{
			hash := os.Args[2]
			if len(os.Args) >= 4 {
				if os.Args[3] == "test" {
					normal = false
				}
			}
			fmt.Println("getting ", hash)
			fmt.Println(string(network.SendLocalGet(hash)))
		}
	case "put":
		{
			data := os.Args[2]
			if len(os.Args) >= 4 {
				if os.Args[3] == "test" {
					normal = false
				}
			}
			fmt.Println("storing ", data)
			fmt.Println(network.SendLocalPut([]byte(data)))
		}

	default:
		printHelpExit("Invalid command.")
	}
	go RestApi()

	//Testing call for Listen

	if normal {
		network.Listen()
	}

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

func IpPortSerialize(myIP net.IP, port int) string {
	localIP := myIP.String() + ":" + strconv.Itoa(port)
	return localIP
}
