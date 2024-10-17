package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
)

// Added for NodeLookup
const alpha = 3 // Number of contacts to retrieve in each round of nodelookup, aka number of parallel queries that concurrently can be sent to other nodes, for speed
const k = 5     // Maximum number of closest contacts to retain in the contactList during lookup.

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

// Added for NodeLookup
type ContactDistance struct {
	Contact          Contact
	DistanceToTarget *KademliaID
	Probed           bool
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

// NodeLookup performs a lookup to find the closest nodes or retrieve data in the network
func (kademlia *Kademlia) NodeLookup(target *Contact, hash string) ([]Contact, Contact, []byte) {

	// Retrieve the closest alpha contacts to the target from the routing table
	closestContacts := kademlia.RoutingTable.FindClosestContacts(target.ID, alpha)

	// Initialize a dynamic list of ContactDistance to hold the list (contactList) of closest contacts
	var contactList []ContactDistance
	for _, contact := range closestContacts {
		// Add each contact to the contactList and sort by distance to the target
		contactList = AddContactToContactList(contactList, contact, target.ID)
	}

	// Keep track of the closest node, aka first item in contactList
	closestNode := contactList[0]

	// Iterative lookup
	for {
		// Retrieve unprobed (not yet queried) contacts from contactList
		unprobedContacts := kademlia.GetNextUnprobedContacts(contactList)

		// If no more unprobed contacts, return gathered contacts
		if len(unprobedContacts) == 0 {
			return GetAllContactsFromContactList(contactList), Contact{}, nil
		}

		// Probe the unprobed contacts and send FindNode messages
		var contactFoundDataOn Contact
		var foundData []byte

		contactList, contactFoundDataOn, foundData = kademlia.SendFindNodeMessagesToUnprobedContacts(contactList, target, hash, unprobedContacts)
		fmt.Println("DEBUG: ContactList after sending messages", contactList)

		// If data is found, return the results along with the contacts
		if foundData != nil {
			fmt.Println("Node lookup complete, data found")
			return GetAllContactsFromContactList(contactList), contactFoundDataOn, foundData
		}

		// Update the closest node from the contact list after probing
		newClosestNode := contactList[0]

		// If the closest node hasn't changed, check for termination conditions
		if closestNode.Contact.ID.Equals(newClosestNode.Contact.ID) {
			// Get the remaining unprobed contacts
			unprobedKClosest := kademlia.GetNextUnprobedContacts(contactList)

			// If enough contacts have been probed or there are no more unprobed contacts, terminate
			if CountProbedInContactList(contactList) >= k || len(unprobedKClosest) == 0 {
				break
			} else {
				// Continue probing the next set of closest unprobed contacts
				newContactList, _, _ := kademlia.SendFindNodeMessagesToUnprobedContacts(contactList, target, hash, unprobedKClosest)
				contactList = newContactList
			}
		} else {
			// Update the closest node and continue the lookup process
			closestNode = newClosestNode
		}
	}

	// Final return: no data found, return all contacts in the list
	fmt.Println("Node lookup complete")
	return GetAllContactsFromContactList(contactList), Contact{}, nil
}

// Add new contact to constactList, avoiding duplicates and sorting the list.
func AddContactToContactList(contactList []ContactDistance, newContact Contact, targetID *KademliaID) []ContactDistance {
	// Check if the contact already is in contactList, if so keep list unchanged
	for _, item := range contactList {
		if item.Contact.ID.Equals(newContact.ID) {
			return contactList
		}
	}

	// Calculate distance between new contact and target
	newDistance := newContact.ID.CalcDistance(targetID)
	newItem := ContactDistance{
		Contact:          newContact,
		DistanceToTarget: newDistance,
		Probed:           false,
	}

	// Add new contact to contactList
	contactList = append(contactList, newItem)

	// Sort contactList by distance to target
	sort.Slice(contactList, func(i, j int) bool {
		return contactList[i].DistanceToTarget.Less(contactList[j].DistanceToTarget)
	})

	// If exceeding k contacts, shorten contactList to k closest contacts
	if len(contactList) > k && k > 0 {
		contactList = contactList[:k]
	}

	return contactList
}

// GetNextUnprobedContacts returns up to alpha unprobed contacts from the contactList, excluding the local node
func (kademlia *Kademlia) GetNextUnprobedContacts(contactList []ContactDistance) []ContactDistance {
	unprobedContacts := make([]ContactDistance, 0, alpha) // Create a slice with up to alpha number of elements

	// Collect unprobed contacts, skip local node
	for _, item := range contactList {
		if !item.Probed && !item.Contact.ID.Equals(kademlia.RoutingTable.me.ID) {
			unprobedContacts = append(unprobedContacts, item)
		}
		// Stop collecting once number of alpha unprobed contacts are reached
		if len(unprobedContacts) == alpha {
			break
		}
	}

	return unprobedContacts
}

// Return all contacts from contactList
func GetAllContactsFromContactList(contactList []ContactDistance) []Contact {
	contacts := make([]Contact, len(contactList)) // Preallocate slice with same number of contacts from the contactList
	for i, item := range contactList {
		contacts[i] = item.Contact // Assign each contact to the preallocated slice
	}
	return contacts
}

// Sends FIND_NODE messages to unprobed contacts in contactList and handles responses.
func (kademlia *Kademlia) SendFindNodeMessagesToUnprobedContacts(
	contactList []ContactDistance, target *Contact, hash string, unprobedContacts []ContactDistance,
) ([]ContactDistance, Contact, []byte) {

	// Initialize channels to collect responses concurrently from up to alpha unprobed contacts
	contactResponses := make(chan Contact, alpha)
	dataResponses := make(chan []byte, alpha)
	contactWithData := make(chan Contact, alpha)

	// Concurrently send FIND_NODE messages to unprobed contacts
	go kademlia.sendFindNodeQueries(unprobedContacts, target, hash, contactResponses, dataResponses, contactWithData)

	// Handle responses, closing channels once all messages have been processed
	defer close(contactResponses)
	defer close(dataResponses)
	defer close(contactWithData)

	fmt.Println("DEBUG: Completed contact probing")

	// Check if any data was returned
	foundContact, foundData := kademlia.handleFoundData(dataResponses, contactWithData)
	if foundData != nil {
		return contactList, foundContact, foundData
	}

	// Process and update the contact list with newly received contacts
	updatedContactList := kademlia.updateContactListWithReceivedContacts(contactList, contactResponses)

	// Mark probed contacts to avoid querying them again
	updatedContactList = kademlia.markContactsAsProbed(updatedContactList, unprobedContacts)

	return updatedContactList, Contact{}, nil
}

// Updates the contactList with the received contacts.
func (kademlia *Kademlia) updateContactListWithReceivedContacts(contactList []ContactDistance, target *Contact, contactsChan chan Contact) []ContactDistance {
	for contact := range contactsChan {
		found := false

		// Check if the contact already exists in the contactList and update it if found
		for i := range contactList {
			if contactList[i].Contact.ID.Equals(contact.ID) {
				contactList[i].Contact = contact
				found = true
				break
			}
		}

		// If the contact was not found, add it to the contactList
		if !found {
			contactList = AddContactToContactList(contactList, contact, target.ID)
		}
	}
	return contactList
}

// Marks the contacts in contactList as probed if they were previously unprobed.
func markContactsAsProbed(contactList []ContactDistance, unprobedContacts []ContactDistance) []ContactDistance {
	// Create a map of unprobed contacts for quick lookup
	unprobedMap := make(map[string]bool)
	for _, contact := range unprobedContacts {
		unprobedMap[contact.Contact.ID.String()] = true
	}

	// Mark contacts in contactList as probed if found in the unprobed map
	for i, contactItem := range contactList {
		if unprobedMap[contactItem.Contact.ID.String()] {
			contactList[i].Probed = true
		}
	}

	return contactList
}

// Counts how many contacts in contactList have been probed.
func CountProbedInContactList(contactList []ContactDistance) int {
	probedCount := 0

	// Loop through contactList and increment the counter for each probed contact
	for _, contact := range contactList {
		if contact.Probed {
			probedCount++
		}
	}

	return probedCount
}

// Sends FIND_NODE or FIND_DATA queries at the same time to unprobed contacts.
func (kademlia *Kademlia) sendFindNodeQueries(unprobedContacts []ContactDistance, target *Contact, hash string, contactResponses chan Contact, dataResponses chan []byte, contactWithDataChan chan Contact) {
	var wg sync.WaitGroup

	queryType := "FIND_NODE"
	if hash != "" {
		queryType = "FIND_DATA"
	}

	// Predefine query function based on whether we are looking for data or contacts
	queryFunc := func(contact Contact) {
		defer wg.Done()
		if queryType == "FIND_NODE" {
			kademlia.findContact(contact, target, contactResponses, dataResponses, contactWithDataChan)
		} else {
			kademlia.findData(contact, hash, dataResponses, contactWithDataChan)
		}
	}

	// Process each unprobed contact concurrently
	for _, contactItem := range unprobedContacts {
		wg.Add(1)
		go queryFunc(contactItem.Contact)
	}

	// Wait for all queries to complete
	wg.Wait()
}

// findContact sends a FIND_NODE message and sends found contacts via channels.
func (kademlia *Kademlia) findContact(contact Contact, target *Contact, contactsChan chan Contact, dataChan chan []byte, contactWithData chan Contact) {
	if contacts, err := kademlia.Network.SendFindContactMessage(&kademlia.RoutingTable.me, &contact, target); err == nil { //MAKE A SendFindContactMessage IN NETWORK??????????????????????
		for _, foundContact := range contacts {
			select {
			case contactsChan <- foundContact:
				dataChan <- nil
				contactWithData <- Contact{}
			default:
				fmt.Printf("DEBUG: Channel full, contact %s not sent\n", foundContact.ID.String())
			}
		}
	} else {
		fmt.Printf("Error sending FIND_NODE to %s: %v\n", contact.ID.String(), err)
	}
}

// findData sends a FIND_DATA message and sends the found data via channels.
func (kademlia *Kademlia) findData(contact Contact, hash string, dataChan chan []byte, contactWithDataChan chan Contact) {
	if _, data, err := kademlia.Network.SendFindDataMessage(&kademlia.RoutingTable.me, &contact, hash); err == nil && data != nil { ////MAKE A SendFindDataMessage IN NETWORK??????????????????????
		fmt.Printf("DEBUG: Found data for hash %s on contact %s\n", hash, contact.ID.String())
		dataChan <- data
		contactWithDataChan <- contact
	} else if err != nil {
		fmt.Printf("Error sending FIND_DATA to %s: %v\n", contact.ID.String(), err)
	}
}

// Checks for data from the data channel and returns the contact and data if available.
func handleFoundData(dataChan chan []byte, contactChanFoundDataOn chan Contact) (Contact, []byte) {
	if data := <-dataChan; data != nil {
		if foundContact := <-contactChanFoundDataOn; foundContact.ID != nil {
			return foundContact, data
		}
	}
	return Contact{}, nil
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
