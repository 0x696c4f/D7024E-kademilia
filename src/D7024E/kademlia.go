package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

// This Kademlia node
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

			if shortList.Len() < network.Node.Alpha {
				contactsToContact = shortList.GetContacts(shortList.Len())

				iterations := shortList.Len() //This is needed, because the shortlist is changed within the function
				for i := 0; i < iterations; i++ {
					//want to pick contact alpha amount of contacts but I don't know if is should be .node .routingtable or other
					network.SendFindContactMessage(&contactsToContact[i], target)
					alreadyContacted = append(alreadyContacted, contactsToContact[i])

				}
			} else {
				contactsToContact = shortList.GetContacts(network.Node.Alpha)
				for i := 0; i < network.Node.Alpha; i++ {
					//want to pick contact alpha amount of contacts but   I don't know if is should be .node .routingtable or other
					network.SendFindContactMessage(&contactsToContact[i], target)
					alreadyContacted = append(alreadyContacted, contactsToContact[i])
				}
			}

			network.Node.ManageShortList(&shortList)

			if shortList.contacts[0].Less(closestContact) { //if new is cloesst
				closestContact = &shortList.contacts[0] //change the new to the accuall closet
			} else { //if the old is closest
				findingClosestNode = false

				lastContactContacts := XorContactLists(shortList.contacts, alreadyContacted)

				for i := 0; i < len(lastContactContacts); i++ {
					network.SendFindContactMessage(&lastContactContacts[i], target)
				}

				network.Node.ManageShortList(&shortList)
			}

		}

		return &shortList
	} else {
		emptyStruct := NewEmptyContactCandidates()
		return &emptyStruct
	}
}

/*
The FIND_VALUE RPC behaves like FIND_NODE, returning the k nodes closest to the target identifier with one exception

	– if the RPC recipient has received a STORE for the given key, it returns the stored value.
*/
func (network *Network) LookupData(hash string) []byte {
	//create a new hashed contact
	hashKademliaID := NewKademliaID(hash)
	hashContact := NewContact(hashKademliaID, "")

	//Find the closest nodes for the key
	shortlist := network.LookupContact(&hashContact)
	var data []byte
	//Send the store RPC
	for _, contact := range shortlist.contacts {
		data = network.SendFindDataMessage(hashKademliaID.String(), contact)
		if data != nil {
			break
		}
	}
	return data
}

// The sender of the STORE RPC provides a key and a block of data and requires that the recipient store the data and make it available for later retrieval by that key.
func (network *Network) Store(data []byte) *KademliaID {
	//create a new hashed contact
	hashInput := HashData(string(data))
	hashKademliaID := NewKademliaID(hashInput)
	hashContact := NewContact(hashKademliaID, "")

	//Find the closest nodes for the key
	closestNodes := network.LookupContact(&hashContact)

	//Send the store RPC
	for _, storeAtNode := range closestNodes.contacts {
		network.SendStoreMessage(data, &storeAtNode)
	}
	storedAt := make([]Contact, len(closestNodes.contacts))
	copy(storedAt, closestNodes.contacts)
	network.Mu.Lock()
	defer network.Mu.Unlock()
	network.Refresh[hashInput] = storedAt
	fmt.Println("Store: ", hashKademliaID)
	return hashKademliaID
}

func (network *Network) Forget(hash string) {
	network.Mu.Lock()
	defer network.Mu.Unlock()
	delete(network.Refresh, hash)
}

// (https://stackoverflow.com/questions/10701874/generating-the-sha-hash-of-a-string-using-golang)
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
