package D7024E

import "fmt"

type Network struct {
}

func Listen(ip string, port int) {
	// TODO
}

func testOtherFile() {
	fmt.Println("pinging")
}

func (network *Network) SendPingMessage(contact *Contact) {

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
