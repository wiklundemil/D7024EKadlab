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

func (kademlia *Kademlia) Store(data []byte){
	
	hashedData := sha1.Sum(data) //This is required also the get it up to the amount of chars to later pass to NewKademliaID due to that function require length of at least 20 chars.
	key := hex.EncodeToString(hashedData[:])

	if _, found := (*kademlia.Data)[key]; found { //If there is data on this key already we return error.
		fmt.Println("Data could not be modified, already exist data with this key.")
	}
	//data, found := kademlia.Data[key] //data = the data associated with key and found is a bool if data exists.
	(*kademlia.Data)[key] = data //We save the data in the correct format, i.e in the map
	fmt.Println(key)
	storeID := NewKademliaID(key)
	contacts := kademlia.LookupContact(storeID)
	
	fmt.Println("Stored data with: ", key)
	
	for _, contact := range contacts {
		fmt.Println("Stored data at:", &contact, " with key: ", key)
		go kademlia.Network.SendStoreMessage(data, &contact)
	}
}
