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

func NewKademlia(ipAddress string) (node Kademlia) {
	ID := NewKademliaID(HashData(ipAddress))
	node.Me = NewContact(ID, ipAddress)
	node.RoutingTable = NewRoutingTable(node.Me)
	node.Alpha = 4
	return
}

func (network *Network) LookupContact(target *Contact) *ContactCandidates {
	//find the closest current contact to the looked upon contact
	closestContactsList := network.Node.RoutingTable.FindClosestContacts(target.ID, network.Node.Alpha)
	fmt.Println("closestContactList - ", closestContactsList)
	//if we have any to go throughsu
	if len(closestContactsList) != 0 {
		//current closest node
		shortList := NewContactCandidates(closestContactsList)
		network.Node.Shortlist = &shortList
		closestContact := &shortList.contacts[0]

		var alreadyContacted []Contact

		findingClosestNode := true
		for findingClosestNode {
			var contactsToContact []Contact
			fmt.Println("New loop")
			if shortList.Len() < network.Node.Alpha {
				contactsToContact = shortList.GetContacts(shortList.Len())
				fmt.Println("contactContacts - ", contactsToContact)

				iterations := shortList.Len() //This is needed, because the shortlist is changed within the function
				for i := 0; i < iterations; i++ {
					//want to pick contact alpha amount of contacts but I don't know if is should be .node .routingtable or other
					network.SendFindContactMessage(&contactsToContact[i], target) //TODO make it a go
					alreadyContacted = append(alreadyContacted, contactsToContact[i])
					//TODO add the once we have contacted to cotacted list
				}
			} else {
				contactsToContact = shortList.GetContacts(network.Node.Alpha)
				for i := 0; i < network.Node.Alpha; i++ {
					//want to pick contact alpha amount of contacts but   I don't know if is should be .node .routingtable or other
					network.SendFindContactMessage(&contactsToContact[i], target) //TODO make it a go
					alreadyContacted = append(alreadyContacted, contactsToContact[i])
				}
			}

			network.Node.ManageShortList(&shortList)
			fmt.Println("managed shortlist-", network.Node.Shortlist)

			if shortList.contacts[0].Less(closestContact) { //if new is cloesst
				closestContact = &shortList.contacts[0] //change the new to the accuall closet
			} else { //if the old is closest
				fmt.Println("ending, check list")
				findingClosestNode = false

				lastContactContacts := XorContactLists(shortList.contacts, alreadyContacted)

				for i := 0; i < len(lastContactContacts); i++ {
					network.SendFindContactMessage(&lastContactContacts[i], target)
				}

				network.Node.ManageShortList(&shortList)
			}

		}
		fmt.Println("Finishing Shortlist")
		fmt.Println(shortList)

		return &shortList
	} else {
		emptyStuct := NewEmptyContactCandidates()
		return &emptyStuct
	}
}

func (kademlia *Kademlia) LookupData(hash string) (data []byte){
	fmt.Println("Looking up",hash)
	// TODO
	return ([]byte("example data"))
}

func (kademlia *Kademlia) Store(data []byte) (hash *KademliaID){
	fmt.Println("Storing",data)
	// TODO
	return NewRandomKademliaID()
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
	shortlist.RemoveContact(&kademlia.RoutingTable.me) //remove ourself
	shortlist.RemoveDublicate()
	if bucketSize < shortlist.Len() {
		shortlist.contacts = shortlist.GetContacts(bucketSize)
	}
}
