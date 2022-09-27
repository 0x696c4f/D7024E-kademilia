package main

import (
	"fmt"
	"testing"
)

func TestContact(t *testing.T) {
	fmt.Println("implement contact testing")

	//Create contact
	contactA := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	contactB := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002")
	contactC := NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002")

	//check distance
	contactA.CalcDistance(contactB.ID)
	contactC.CalcDistance(contactB.ID)

	//check less and string
	if contactA.Less(&contactC) {
		fmt.Println("Contact A (", contactA.String(), ") is closer to Contact B (", contactB.String())
	} else {
		fmt.Println("Contact C (", contactB.String(), ") is closer to Contact B(", contactB.String())
	}
	if contactC.Less(&contactA) {
		fmt.Println("Contact A (", contactA.String(), ") is closer to Contact B (", contactB.String())
	} else {
		fmt.Println("Contact C (", contactB.String(), ") is closer to Contact B(", contactB.String())
	}

	//NewContactCandidate
	contactCandidateA := NewEmptyContactCandidates()

	fmt.Println(contactCandidateA)
	//Add item list to it with items
	var contactList []Contact
	contactList = append(contactList, contactA)
	contactList = append(contactList, contactB)
	contactList = append(contactList, contactC)
	fmt.Println(contactList)

	contactCandidateA.Append(contactList)

	//Create new ContactCandidates with a already made list
	contactCandidateB := NewContactCandidates(contactList)
	fmt.Println(contactCandidateB)

	//get the first two thirds of the list
	length := contactCandidateA.Len()
	twoThirdsContactList := contactCandidateA.GetContacts((length / 3) * 2)
	fmt.Println("Delete a third of the items", twoThirdsContactList)

	contactCandidateA.Swap(0, 1)
	contactCandidateA.Swap(-3030, -3030)
	contactCandidateA.Swap(3030, 3030)
	contactCandidateA.Swap(-3030, 3030)
	contactCandidateA.Swap(3030, -3030)
	contactCandidateA.Swap(1, 3030)
	contactCandidateA.Swap(1, -3030)
	contactCandidateA.Swap(3030, 1)
	contactCandidateA.Swap(-3030, 1)

	//remove a contact from the list
	fmt.Println("Before removal of contact = ", contactCandidateA)
	contactCandidateA.RemoveContact(&contactA)
	fmt.Println("After removal of contact = ", contactCandidateA)

	//remove a dublicate
	contactCandidateA.contacts = append(contactCandidateA.contacts, contactB)
	fmt.Println("Added dublicate  = ", contactCandidateA)
	contactCandidateA.RemoveDublicate()
	fmt.Println("remove dublicate = ", contactCandidateA)

	// return the item online in one
	fmt.Println("Before xor A = ", contactCandidateA)
	fmt.Println("Before xor B = ", contactCandidateB)
	xorContactList := XorContactLists(contactCandidateB.contacts, contactCandidateA.contacts)
	fmt.Println("result = ", xorContactList)

	fmt.Println("-------------------------")
	fmt.Println("")
}
