package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

// This Kademlia node
type Kademlia struct {
	routingTable *RoutingTable //Everyone else infromation
	network      *Network
}

const alpha = 3

func NewKademlia(ipAddress string) (node Kademlia) {
	ID := NewKademliaID(HashData(ipAddress))
	node.routingTable = NewRoutingTable(NewContact(ID, ipAddress))

	return
}

/*
FIND_NODE
Note that in the Kademlia paper, objects are referred to as values and their hashes as keys.
*/
func (kademlia *Kademlia) LookupContact(target *Contact) *ContactCandidates {
	fmt.Println("Looking for contact: " + target.ID.String())

	//Find k closest nodes
	closestNodes := kademlia.routingTable.FindClosestContacts(target.ID, alpha)
	if len(closestNodes) != 0 {
		//initiate closeNode
		closestNodesToContact := ContactCandidates{closestNodes}
		closestNodesToContact.Sort()
		closestContact := &closestNodesToContact.contacts[0]
		//initiate candidateList
		var candidateList ContactCandidates
		candidateList.contacts = closestNodes

		contactedList := ContactCandidates{}

		closerNodeFound := true
		for closerNodeFound {
			//Send FIND_NODE rpc to alpha number of contacts in the candidateList
			numCandidates := candidateList.GetContacts(alpha)
			for _, contact := range numCandidates {
				go kademlia.network.SendFindContactMessage(&contact, target) //using function in network.go (line 139-152)
				contactedList.Append([]Contact{contact})
			}
			kademlia.manageCandidateList(len(numCandidates), &candidateList)

			if candidateList.contacts[0].Less(closestContact) {
				closestContact = &candidateList.contacts[0]
			} else {
				closerNodeFound = false
				//Find closest nodes that have not yet been contacted
				nodesToContact := findNotContactedNodes(&candidateList, &contactedList)
				nodesToContact.Sort()
				nodesToContact.RemoveDuplicates()
				//Send an rpc to each of the k closest modes that has not already been contacted
				for _, nodeToContact := range nodesToContact.contacts {
					go kademlia.network.SendFindContactMessage(&nodeToContact, target)
				}
				kademlia.manageCandidateList(nodesToContact.Len(), &candidateList)
				//Remove all inactive nodes from the candidateList
				candidateList.contacts = removeInactiveNodes(candidateList, kademlia.network.inactiveNodes) //using inactiveNodes from network struct
			}
		}
		return &candidateList
	} else {
		return &ContactCandidates{[]Contact{}}
	}
}

// Manages the candidateList by sorting, removing duplicates, removing "itself" and keeps the candidateList of k size
func (kademlia *Kademlia) manageCandidateList(num int, candidateList *ContactCandidates) {
	for i := 0; i < num; i++ {
		newCandidateList := <-kademlia.network.candidateListCh //using candidateListCh from network struct
		candidateList.Append(newCandidateList)
		candidateList.Sort()
		candidateList.RemoveDuplicates()
		candidateList.RemoveContact(&kademlia.routingTable.me) //remove self from candidateList
		candidateList.contacts = candidateList.GetContacts(bucketSize)
	}
}

// Returns the contacts in the candidateList that have not been contacted
func findNotContactedNodes(candidateList *ContactCandidates, contactedNodes *ContactCandidates) ContactCandidates {
	notContactedList := make([]Contact, 0)
	for _, contact := range candidateList.contacts {
		notContacted := true
		for _, contactedNode := range contactedNodes.contacts {
			if contact.ID == contactedNode.ID {
				notContacted = false
			}
		}
		if notContacted {
			notContactedList = append(notContactedList, contact)
		}
	}
	return ContactCandidates{notContactedList}
}

// Returns a clean candidateList without the inactive nodes
func removeInactiveNodes(candidateList ContactCandidates, inactiveNodes ContactCandidates) []Contact {
	cleanCandidateList := make([]Contact, 0)
	for _, contact := range candidateList.contacts {
		isActive := true
		for _, inactiveNode := range inactiveNodes.contacts {
			if contact.ID == inactiveNode.ID {
				isActive = false
			}
		}
		if isActive {
			cleanCandidateList = append(cleanCandidateList, contact)
		}
	}
	return cleanCandidateList
}

/*
FIND_VALUE
*/
func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

// (https://stackoverflow.com/questions/10701874/generating-the-sha-hash-of-a-string-using-golang)
func HashData(msg string) string {
	hash := sha1.New()
	hash.Write([]byte(msg))
	hashString := hex.EncodeToString(hash.Sum(nil))
	return hashString
}
