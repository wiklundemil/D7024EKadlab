package kademlia

import (
	"crypto/sha1"
	"fmt"
	"testing"
)

func TestKademliaLookup(t *testing.T) {
	me := NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000")
	network := &Network{
		responseChan: make(chan Response),
		connection:         nil,
	}
	kademlia := &Kademlia{
		RoutingTable: NewRoutingTable(me),
		Network:      network,
	}

	contact1 := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8001")
	kademlia.RoutingTable.AddContact(contact1)

	target := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "")
	closestContacts := kademlia.LookupContact(&target)

	if len(closestContacts) == 0 {
		t.Fatalf("Expected to find at least one contact, but found none")
	}

	if !closestContacts[0].ID.Equals(contact1.ID) {
		t.Fatalf("Expected first contact to be contact1, got %v", closestContacts[0].ID.String())
	}
}

func TestStoreAndRetrieveData(t *testing.T) {
	me := NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000")
	network := &Network{
		responseChan: make(chan Response),
		connection:         nil,
	}
	kademlia := &Kademlia{
		RoutingTable: NewRoutingTable(me),
		Network:      network,
		Data:         &map[string][]byte{},
	}

	data := []byte("Hello, Kademlia!")
	hash := kademlia.GenerateHash(data)
	kademlia.Store(hash, data)

	retrievedData, err := kademlia.Get(hash)
	if err != nil || retrievedData == nil {
		t.Fatalf("Failed to retrieve stored data: %v", err)
	}

	if string(retrievedData) != string(data) {
		t.Fatalf("Retrieved data does not match stored data. Expected: %s, got: %s", data, retrievedData)
	}
}

func (kademlia *Kademlia) GenerateHash(data []byte) string {
	hash := sha1.Sum(data)
	return fmt.Sprintf("%x", hash)
}

func (kademlia *Kademlia) Get(hash string) ([]byte, error) {
	if data, exists := (*kademlia.Data)[hash]; exists {
		return data, nil
	}
	return nil, fmt.Errorf("data not found for hash: %s", hash)
}
