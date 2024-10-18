package kademlia

import (
	"crypto/sha1"
	"fmt"
	"testing"
)

func TestKademliaLookup(t *testing.T) {
	me := NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000")
	network := &Network{
		Self:         &me,
		RoutingTable: NewRoutingTable(me),
	}
	kademlia := &Kademlia{
		RoutingTable: network.RoutingTable,
		Network:      network,
	}

	// Add contacts to the routing table
	contact1 := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8001")
	network.RoutingTable.AddContact(contact1)

	// Perform lookup for a contact
	target := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "")
	closestContacts, err := kademlia.LookupContact(&target)
	if err != nil {
		t.Fatalf("Lookup failed: %v", err)
	}

	if len(closestContacts) == 0 {
		t.Fatalf("Expected to find at least one contact, but found none")
	}

	if !closestContacts[0].ID.Equals(contact1.ID) {
		t.Fatalf("Expected first contact to be contact1")
	}
}

func TestStoreAndRetrieveData(t *testing.T) {
	me := NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000")
	network := &Network{
		Self:         &me,
		RoutingTable: NewRoutingTable(me),
	}
	kademlia := &Kademlia{
		RoutingTable: network.RoutingTable,
		Network:      network,
		Data:         &map[string][]byte{},
	}

	// Store data
	data := []byte("Hello, Kademlia!")
	hash := kademlia.GenerateHash(data) // Assuming GenerateHash method is now implemented
	kademlia.Store(data)

	// Verify if data can be retrieved
	retrievedData, _, err := kademlia.Get(hash)
	if err != nil || retrievedData == nil {
		t.Fatalf("Failed to retrieve stored data")
	}

	if string(retrievedData) != string(data) {
		t.Fatalf("Retrieved data does not match stored data")
	}
}

// Add GenerateHash method to the Kademlia struct
func (kademlia *Kademlia) GenerateHash(data []byte) string {
	hash := sha1.Sum(data)
	return fmt.Sprintf("%x", hash)
}
