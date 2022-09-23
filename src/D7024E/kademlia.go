package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

//This Kademlia node
type Kademlia struct {
	RoutingTable *RoutingTable //Everyone else infromation
	Alpha        int
	Network      *Network
}

func (network *Network) NewKademlia(ipAddress string) (node Kademlia) {
	ID := NewKademliaID(HashData(ipAddress))
	node.RoutingTable = NewRoutingTable(NewContact(ID, ipAddress))
	node.Network = network
	node.Alpha = 4
	return
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	fmt.Println("hejsan")
	//find the closest current contact to the looked upon contact
	fmt.Println(kademlia.Alpha)
	closestContactsList := kademlia.RoutingTable.FindClosestContacts(target.ID, kademlia.Alpha)
	fmt.Println(closestContactsList)
	bucketIndex := kademlia.RoutingTable.getBucketIndex(target.ID)
	bucket := kademlia.RoutingTable.buckets[bucketIndex].list
	fmt.Println(bucket.Len())

	//if we have any to go throughsu
	if len(closestContactsList) != 0 {
		fmt.Println("komm")
		//current closest node
		contactCandidate := NewContactCandidates(closestContactsList)
		contactCandidate.Sort()
		closestContact := &contactCandidate.contacts[0]

		//keep track of what nodes we have gone through
		//contactedList := NewEmptyContactCandidates()

		//keep track of the nodes we will go through, first shortlist sorted
		shortList := NewContactCandidates(closestContactsList)

		//contact the list to se if you have any closer to what you are looking for
		closestNode := true
		fmt.Println("before for")
		for closestNode {
			fmt.Println("afterfor")
			var contactContacts []Contact
			if shortList.Len() < kademlia.Alpha {
				contactContacts = shortList.GetContacts(shortList.Len())
				for i := 0; i < shortList.Len(); i++ {
					//want to pick contact alpha amount of contacts but I don't know if is should be .node .routingtable or other
					kademlia.Network.SendFindContactMessage(&contactContacts[i], target) //TODO make it a go
					//TODO add the once we have contacted to cotacted list
				}
			} else {
				contactContacts = shortList.GetContacts(kademlia.Alpha)
				for i := 0; i < kademlia.Alpha; i++ {
					//want to pick contact alpha amount of contacts but I don't know if is should be .node .routingtable or other
					kademlia.Network.SendFindContactMessage(&contactContacts[i], target) //TODO make it a go
					//TODO add the once we have contacted to cotacted list
				}
			}

			//TODO manage shortlist with the new data from called nodes

			if shortList.contacts[0].Less(closestContact) { //if new is cloesst
				closestContact = &shortList.contacts[0] //change the new to the accuall closet
			} else { //if the old is closest
				//TODO break loop
				/*
					send a new sendfindcontactmessage to alpha contact not contacted yet.
				*/

			}
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
