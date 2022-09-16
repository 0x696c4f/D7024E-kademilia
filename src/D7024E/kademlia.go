package main

import (
	"crypto/sha1"
	"encoding/hex"
)

//This Kademlia node
type Kademlia struct {
	routingTable *RoutingTable //Everyone else infromation
}

func NewKademlia(ipAddress string) (node Kademlia) {
	ID := NewKademliaID(HashData(ipAddress))
	node.routingTable = NewRoutingTable(NewContact(ID, ipAddress))

	return
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

//(https://stackoverflow.com/questions/10701874/generating-the-sha-hash-of-a-string-using-golang)
func HashData(msg string) string {
	hash := sha1.New()
	hash.Write([]byte(msg))
	hashString := hex.EncodeToString(hash.Sum(nil))
	return hashString
}
