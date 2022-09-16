package main

import (
	"container/list"
	"fmt"
)

const bucketSize = 20

// RoutingTable definition
// keeps a refrence contact of me and an array of buckets
type RoutingTable struct {
	me      Contact //my own information
	buckets [IDLength * 8]*bucket
}

// NewRoutingTable returns a new instance of a RoutingTable
func NewRoutingTable(me Contact) *RoutingTable {
	routingTable := &RoutingTable{}
	for i := 0; i < IDLength*8; i++ {
		routingTable.buckets[i] = newBucket()
	}
	routingTable.me = me
	return routingTable
}

// AddContact add a new contact to the correct Bucket
func (routingTable *RoutingTable) AddContact(contact Contact) {
	bucketIndex := routingTable.getBucketIndex(contact.ID)
	bucket := routingTable.buckets[bucketIndex]
	bucket.AddContact(contact)
}

func (routingTable *RoutingTable) RemoveContact(contact Contact) {
	bucketIndex := routingTable.getBucketIndex(contact.ID)
	bucket := routingTable.buckets[bucketIndex].list

	var element *list.Element
	for item := bucket.Front(); item != nil; item = item.Next() { //check if contact is within the bucket
		if contact.ID == item.Value.(Contact).ID {
			element = item
		}
	}
	if element != nil {
		bucket.Remove(element)
	}
}

// FindClosestContacts finds the count closest Contacts to the target in the RoutingTable
func (routingTable *RoutingTable) FindClosestContacts(target *KademliaID, count int) []Contact {
	var candidates ContactCandidates
	bucketIndex := routingTable.getBucketIndex(target)
	bucket := routingTable.buckets[bucketIndex]

	candidates.Append(bucket.GetContactAndCalcDistance(target))

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < IDLength*8) && candidates.Len() < count; i++ {
		if bucketIndex-i >= 0 {
			bucket = routingTable.buckets[bucketIndex-i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
		if bucketIndex+i < IDLength*8 {
			bucket = routingTable.buckets[bucketIndex+i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
	}

	candidates.Sort()

	if count > candidates.Len() {
		count = candidates.Len()
	}

	return candidates.GetContacts(count)
}

// getBucketIndex get the correct Bucket index for the KademliaID
func (routingTable *RoutingTable) getBucketIndex(id *KademliaID) int {
	distance := id.CalcDistance(routingTable.me.ID)
	for i := 0; i < IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}

	return IDLength*8 - 1
}

func (network *Network) AddToRoutingTable(contact Contact) {

	//I take this as we go into the right bucket
	bucketIndex := network.node.routingTable.getBucketIndex(contact.ID)
	bucket := network.node.routingTable.buckets[bucketIndex].list

	var element *list.Element
	for item := bucket.Front(); item != nil; item = item.Next() { //check if contact is within the bucket
		if contact.ID == item.Value.(Contact).ID {
			element = item
		}
	}

	if element != nil { //the element is already in the bucket

		bucket.MoveToFront(element)
		fmt.Println("length 2: ", bucket.Len())
		fmt.Println("TODO move to the front")

	} else if element == nil && bucket.Len() < bucketSize { //not in bucket and the bucket is not full
		network.node.routingTable.AddContact(contact)

	} else if element == nil && bucket.Len() >= bucketSize { //not in bucket and the bucket is full

		backElement := bucket.Back()
		backcontact := backElement.Value.(Contact)
		response, _ := network.SendPingMessage(&backcontact)

		if response.RPC != "pong" { //if the last element doesn't responde
			bucket.Remove(backElement)
			network.node.routingTable.AddContact(contact)
		} else {
			bucket.MoveToFront(backElement)
		}

	}
}
