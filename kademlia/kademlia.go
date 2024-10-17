package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

type Kademlia struct {
	RoutingTable *RoutingTable
	Network      *Network

	// Initialize data storage as a map with string keys and byte slice values
	// This map will be used to store data in the Kademlia instance.
	// - Keys (string): These represent unique identifiers for the data entries.
	// - Values ([]byte): These are byte slices, allowing for flexible storage of various types of data, such as text, files, or other serialized objects.
	// Using a map provides efficient lookups, insertions, and deletions, which is crucial for the performance of the distributed system.
	Data *map[string][]byte
}

func (kademlia *Kademlia) LookupContact(target *Contact) ([]Contact, error) {
	closestContacts := kademlia.RoutingTable.FindClosestContacts(target.ID, 10) //finding 10 closest nodes to target.ID from routing table

	if len(closestContacts) == 0 {
		return nil, fmt.Errorf("no contacts found for target ID: %s", target.ID)
	}

	return closestContacts, nil
}

// Checks if data for a hash exists locally and returns data
// If not found locally, a contact based on the hash is created and closesst contacts that may have the data are returned instead
func (kademlia *Kademlia) LookupData(hash string) ([]byte, []Contact, error) {
	if data, ok := (*kademlia.Data)[hash]; ok {
		return data, nil, nil // Return data if found
	}

	contact := NewContact(NewKademliaID(hash), "") // Create a contact

	closestContacts, err := kademlia.LookupContact(&contact) // Find closest contacts
	if err != nil {
		return nil, nil, fmt.Errorf("failed to lookup contacts: %w", err)
	}

	return nil, closestContacts, nil // Return; nil for data not found locally and closest contacts, no error
}

func (kademlia *Kademlia) Store(data []byte) {
	hashedData := sha1.Sum(data)             // Generate SHA-1 hash of the data
	key := hex.EncodeToString(hashedData[:]) // Encode the hash to a string

	if _, found := (*kademlia.Data)[key]; found { // If there is data with this key already, print a message
		fmt.Println("Data could not be modified, already exist data with this key.")
		return
	}

	// Store the data in the map
	(*kademlia.Data)[key] = data
	fmt.Println(key)

	// Create a contact for the store ID
	storeID := NewKademliaID(key)
	contact := NewContact(storeID, "")

	contacts, err := kademlia.LookupContact(&contact)
	if err != nil {
		fmt.Printf("Failed to lookup contacts for storage: %v\n", err)
		return
	}

	fmt.Println("Stored data with key:", key)

	// Send store messages to the closest contacts
	for _, contact := range contacts {
		fmt.Printf("Stored data at contact: %s with key: %s\n", contact.Address, key)
		go kademlia.Network.SendStoreMessage(data, &contact)
	}
}

func (kademlia *Kademlia) Get(hash string) ([]byte, []Contact, error) {
	// Check if the data exists locally in the node’s data store
	data, closestContacts, err := kademlia.LookupData(hash)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to lookup data: %w", err)
	}

	// If the data is found locally, return it
	if data != nil {
		fmt.Println("Data found locally.")
		return data, nil, nil
	}

	// If the data is not found locally, return the closest contacts
	fmt.Println("Data not found locally. Closest contacts returned.")
	return nil, closestContacts, nil
}
