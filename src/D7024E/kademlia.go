package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

//This Kademlia node
type Kademlia struct {
	Me           Contact
	RoutingTable *RoutingTable //Everyone else infromation
	Alpha        int
	Shortlist    *ContactCandidates
}

func (network *Network) NewKademlia(ipAddress string) (node Kademlia) {
	ID := NewKademliaID(HashData(ipAddress))
	node.Me = NewContact(ID, ipAddress)
	node.RoutingTable = NewRoutingTable(node.Me)
	node.Alpha = 4
	return
}

func (network *Network) LookupContact(target *Contact) {
	fmt.Println("hejsan", target)
	//find the closest current contact to the looked upon contact
	fmt.Println(network.Node.Alpha)
	closestContactsList := network.Node.RoutingTable.FindClosestContacts(target.ID, network.Node.Alpha)
	fmt.Println(closestContactsList)
	bucketIndex := network.Node.RoutingTable.getBucketIndex(target.ID)
	bucket := network.Node.RoutingTable.buckets[bucketIndex].list
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
		var shortList ContactCandidates

		shortList.contacts = closestContactsList
		network.Node.Shortlist = &shortList

		shortList.Sort()
		//contact the list to se if you have any closer to what you are looking for
		closestNode := true
		fmt.Print("before for - ")
		for closestNode {
			fmt.Println("afterfor")
			var contactContacts []Contact
			fmt.Println("a-", shortList)

			if shortList.Len() < network.Node.Alpha {
				contactContacts = shortList.GetContacts(shortList.Len())
				iterations := shortList.Len()
				for i := 0; i < iterations; i++ {
					//want to pick contact alpha amount of contacts but I don't know if is should be .node .routingtable or other
					network.SendFindContactMessage(&contactContacts[i], target) //TODO make it a go
					//TODO add the once we have contacted to cotacted list
				}
			} else {
				contactContacts = shortList.GetContacts(network.Node.Alpha)
				for i := 0; i < network.Node.Alpha; i++ {
					//want to pick contact alpha amount of contacts but   I don't know if is should be .node .routingtable or other
					network.SendFindContactMessage(&contactContacts[i], target) //TODO make it a go
					//TODO add the once we have contacted to cotacted list
				}
			}
			fmt.Println("should have long shortlist-", shortList)
			network.Node.ManageShortList(&shortList)

			fmt.Println("left if statement")
			//TODO manage shortlist with the new data from called nodes

			if shortList.contacts[0].Less(closestContact) { //if new is cloesst
				closestContact = &shortList.contacts[0] //change the new to the accuall closet
			} else { //if the old is closest

				fmt.Println("else statement")

				//TODO break loop
				/*
					send a new sendfindcontactmessage to alpha contact not contacted yet.
				*/

			}
			closestNode = false

			fmt.Println("Set false")
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

func (kademlia *Kademlia) ManageShortList(shortlist *ContactCandidates) {
	shortlist.Sort()
	shortlist.RemoveContact(&kademlia.RoutingTable.me)
	shortlist.RemoveDublicate()
	if bucketSize < shortlist.Len() {
		shortlist.contacts = shortlist.GetContacts(bucketSize)
	}
}
