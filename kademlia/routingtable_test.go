package kademlia

import (
	"fmt"
	"testing"
)

func TestRoutingTableAddAndFind(t *testing.T) {
	me := NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000")
	routingTable := NewRoutingTable(me)

	contact1 := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8001")
	routingTable.AddContact(contact1)

	closestContacts := routingTable.FindClosestContacts(contact1.ID, 1)
	if len(closestContacts) == 0 || !closestContacts[0].ID.Equals(contact1.ID) {
		t.Fatalf("Expected to find contact1 in the routing table")
	}
}

func (routingTable *RoutingTable) RemoveContact(contact Contact) {
	bucketIndex := routingTable.getBucketIndex(contact.ID)
	bucket := routingTable.buckets[bucketIndex]

	for e := bucket.list.Front(); e != nil; e = e.Next() {
		if e.Value.(Contact).ID.Equals(contact.ID) {
			bucket.list.Remove(e)
			return
		}
	}
}

func TestRoutingTableFullBucket(t *testing.T) {
	me := NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000")
	routingTable := NewRoutingTable(me)

	// Add more contacts than the bucket size
	for i := 0; i < bucketSize+1; i++ {
		contact := NewContact(NewKademliaID(fmt.Sprintf("%040d", i)), fmt.Sprintf("localhost:%d", 8000+i))
		routingTable.AddContact(contact)
	}

	if routingTable.buckets[0].list.Len() > bucketSize {
		t.Fatalf("Expected bucket size to not exceed %d", bucketSize)
	}
}

func TestRoutingTableRemoveContact(t *testing.T) {
	me := NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000")
	routingTable := NewRoutingTable(me)

	contact := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8001")
	routingTable.AddContact(contact)

	// Remove contact
	routingTable.RemoveContact(contact)

	closestContacts := routingTable.FindClosestContacts(contact.ID, 1)
	if len(closestContacts) > 0 {
		t.Fatalf("Failed to remove contact from routing table")
	}
}
