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
	//fmt.Println("closestContactList - ", closestContactsList)
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
			//fmt.Println("New loop")
			if shortList.Len() < network.Node.Alpha {
				contactsToContact = shortList.GetContacts(shortList.Len())
				//fmt.Println("contactContacts - ", contactsToContact)

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
			//fmt.Println("managed shortlist-", network.Node.Shortlist)

			if shortList.contacts[0].Less(closestContact) { //if new is cloesst
				closestContact = &shortList.contacts[0] //change the new to the accuall closet
			} else { //if the old is closest
				//fmt.Println("ending, check list")
				findingClosestNode = false

				lastContactContacts := XorContactLists(shortList.contacts, alreadyContacted)

				for i := 0; i < len(lastContactContacts); i++ {
					network.SendFindContactMessage(&lastContactContacts[i], target)
				}

				network.Node.ManageShortList(&shortList)
			}

		}
		//fmt.Println("Finishing Shortlist")
		//fmt.Println(shortList)

		return &shortList
	} else {
		emptyStruct := NewEmptyContactCandidates()
		return &emptyStruct
	}
}

/*
The FIND_VALUE RPC behaves like FIND_NODE, returning the k nodes closest to the target identifier with one exception

	– if the RPC recipient has received a STORE for the given key, it returns the stored value.

A FIND_VALUE RPC includes a B=160-bit key. If a corresponding value is present on the recipient, the associated data is returned.
Otherwise the RPC is equivalent to a FIND_NODE and a set of k triples is returned.
This is a primitive operation, not an iterative one.
*/
func (network *Network) LookupData(hash string) []byte { // TODO
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

/*
The sender of the STORE RPC provides a key and a block of data and requires that the recipient store the data and make it available for later retrieval by that key.

While this is not formally specified, it is clear that the initial STORE message must contain in addition to the message ID at least the data to be stored (including its length)
and the associated key.
As the transport may be UDP, the message needs to also contain at least the nodeID of the sender, and the reply the nodeID of the recipient.

The reply to any RPC should also contain an indication of the result of the operation. For example, in a STORE while no maximum data length has been specified,
it is clearly possible that the receiver might not be able to store the data, either because of lack of space or because of an I/O error.

For efficiency, the STORE RPC should be two-phase.
In the first phase the initiator sends a key and possibly length and the recipient replies with either something equivalent to OK or a code signifying that it already
has the value or some other status code.
If the reply was OK, then the initiator may send the value.

Some consideration should also be given to the development of methods for handling hierarchical data.
Some values will be small and will fit in a UDP datagram. But some messages will be very large, over say 5 GB, and will need to be chunked.
The chunks themselves might be very large relative to a UDP packet, typically on the order of 128 KB, so these chunks will have to be shredded into individual UDP packets.
*/
func (network *Network) Store(data []byte) *KademliaID { // TODO
	//create a new hashed contact
	//fmt.Println("Data to be hashed #2: ", data)
	hashInput := HashData(string(data))
	//fmt.Println("Hash #2: ", hashInput)
	hashKademliaID := NewKademliaID(hashInput)
	hashContact := NewContact(hashKademliaID, "")

	//Find the closest nodes for the key
	closestNodes := network.LookupContact(&hashContact)

	//Send the store RPC
	for _, storeAtNode := range closestNodes.contacts {
		network.SendStoreMessage(data, &storeAtNode)
	}
	fmt.Println("Store: ", hashKademliaID)
	return hashKademliaID

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
