package main

import (
	"fmt"
	"sort"
	"sync"
)

// Contact definition
// stores the KademliaID, the ip address and the distance
type Contact struct {
	ID       *KademliaID
	Address  string
	Distance *KademliaID
}

// NewContact returns a new instance of a Contact
func NewContact(id *KademliaID, address string) Contact {
	return Contact{id, address, nil}
}

// CalcDistance calculates the distance to the target and
// fills the contacts distance field
func (contact *Contact) CalcDistance(target *KademliaID) {
	contact.Distance = contact.ID.CalcDistance(target)
}

// Less returns true if contact.distance < otherContact.distance
func (contact *Contact) Less(otherContact *Contact) bool {
	return contact.Distance.Less(otherContact.Distance)
}

// String returns a simple string representation of a Contact
func (contact *Contact) String() string {
	return fmt.Sprintf(`contact("%s", "%s")`, contact.ID, contact.Address)
}

// ContactCandidates definition
// stores an array of Contacts
type ContactCandidates struct {
	sync.Mutex
	contacts []Contact
}

func NewContactCandidates(contactsInput []Contact) ContactCandidates {
	candidates := ContactCandidates{
		contacts: contactsInput,
	}
	return candidates
}

func NewEmptyContactCandidates() ContactCandidates {
	candidates := ContactCandidates{}
	return candidates
}

// Append an array of Contacts to the ContactCandidates
func (candidates *ContactCandidates) Append(contacts []Contact) {
	candidates.contacts = append(candidates.contacts, contacts...)
}

// GetContacts returns the first count number of Contacts
func (candidates *ContactCandidates) GetContacts(count int) []Contact {
	return candidates.contacts[:count]
}

// Sort the Contacts in ContactCandidates
func (candidates *ContactCandidates) Sort() {
	sort.Sort(candidates)
}

// Len returns the length of the ContactCandidates
func (candidates *ContactCandidates) Len() int {
	return len(candidates.contacts)
}

// Swap the position of the Contacts at i and j
// WARNING does not check if either i or j is within range
func (candidates *ContactCandidates) Swap(i, j int) {
	candidates.contacts[i], candidates.contacts[j] = candidates.contacts[j], candidates.contacts[i]
}

// Less returns true if the Contact at index i is smaller than
// the Contact at index j
func (candidates *ContactCandidates) Less(i, j int) bool {
	return candidates.contacts[i].Less(&candidates.contacts[j])
}

func XorContactLists(list1 []Contact, list2 []Contact) []Contact {
	sumList := make([]Contact, 0)

	for i := 0; i < len(list1); i++ {
		place := true
		for j := 0; j < len(list2); j++ {
			if list1[i].ID.Equals(list2[j].ID) {
				place = false
				break
			}
		}
		if place {
			sumList = append(sumList, list1[i])
		}
	}

	return sumList
}

func (candidates *ContactCandidates) RemoveContact(removeTheContact *Contact) {
	tempContactList := make([]Contact, 0)
	for i := 0; i < candidates.Len(); i++ {
		if !removeTheContact.ID.Equals(candidates.contacts[i].ID) {
			tempContactList = append(tempContactList, candidates.contacts[i])
		}
	}
	candidates.contacts = tempContactList
}

func (candidates *ContactCandidates) RemoveDublicate() {
	cleanList := make([]Contact, 0)
	for i := 0; i < candidates.Len(); i++ {
		place := true
		for j := 0; j < len(cleanList); j++ {
			if cleanList[j].ID.Equals(candidates.contacts[i].ID) {
				place = false
				//break
			}
		}
		if place {
			cleanList = append(cleanList, candidates.contacts[i])
		}
	}

	candidates.contacts = cleanList
}
