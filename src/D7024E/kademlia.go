package main

import (
	"crypto/sha1"
	"encoding/hex"
)

//This Kademlia node
type Kademlia struct {
	routingTable *RoutingTable //Everyone else infromation
	alpha        int
}

func NewKademlia(ipAddress string) (node Kademlia) {
	ID := NewKademliaID(HashData(ipAddress))
	node.routingTable = NewRoutingTable(NewContact(ID, ipAddress))

	return
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	//find the closest current contact to the looked upon contact
	closestContactsList := kademlia.routingTable.FindClosestContacts(target.ID, kademlia.alpha)

	//if we have any to go through
	if len(closestContactsList) != 0 {
		//current closest node
		contactCandidate := NewContactCandidates(closestContactsList)
		closestContact := &contactCandidate.Sort().contacts[0]

		//keep track of what nodes we have gone through
		contactedList := NewEmptyContactCandidates()

		//keep track of the nodes we will go through, first shortlist
		shortList := NewContactCandidates(closestContactsList)

		//contact the list to se if you have any closer to what you are looking for
		closestNode := true
		for closestNode {
			closestNode = false
		}
	}

	/* The goal is to find a specified contact. How?
	1: create a newtwork and add us the kademlia in to it.
	2: Get the closest nodes within the routingTable
	3: pick out the closest node out of the closest nodes
	4: create for loop to find closer nodes
		4.1: store which has been contacted
		4.2: store the currently known closest nodes
		4:3: store the currently closest node
		--------------------------------------
		5: send contact message to alpha number of items
			5.1: add to contacted nodes
			5.2 delete from closest nodes list
		6: Check if we found the closest contact
			6.1: if not
			6.2: If yes break loop
	*/
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
