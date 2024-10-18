package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
)

// Added for NodeLookup
const alpha = 3 // Number of contacts to retrieve in each round of node lookup, aka number of parallel queries that can be sent to other nodes for speed
const k = 5     // Maximum number of closest contacts to retain in the contactList during lookup.

type Kademlia struct {
	RoutingTable *RoutingTable
	Network      *Network

	// Data storage map
	Data *map[string][]byte
}

// Added for NodeLookup
type ContactDistance struct {
	Contact          Contact
	DistanceToTarget *KademliaID
	Probed           bool
}

func (kademlia *Kademlia) LookupContact(target *Contact) ([]Contact, error) {
	closestContacts := kademlia.RoutingTable.FindClosestContacts(target.ID, 10) // finding 10 closest nodes to target.ID from routing table

	if len(closestContacts) == 0 {
		return nil, fmt.Errorf("no contacts found for target ID: %s", target.ID)
	}

	return closestContacts, nil
}

func (kademlia *Kademlia) LookupData(hash string) ([]byte, []Contact, error) {
	if data, ok := (*kademlia.Data)[hash]; ok {
		return data, nil, nil // Return data if found
	}

	contact := NewContact(NewKademliaID(hash), "")           // Create a contact
	closestContacts, err := kademlia.LookupContact(&contact) // Find closest contacts
	if err != nil {
		return nil, nil, fmt.Errorf("failed to lookup contacts: %w", err)
	}

	return nil, closestContacts, nil // Return nil for data not found locally and closest contacts, no error
}

// NodeLookup performs a lookup to find the closest nodes or retrieve data in the network
func (kademlia *Kademlia) NodeLookup(target *Contact, hash string) ([]Contact, Contact, []byte) {
	closestContacts := kademlia.RoutingTable.FindClosestContacts(target.ID, alpha)

	var contactList []ContactDistance
	for _, contact := range closestContacts {
		contactList = AddContactToContactList(contactList, contact, target.ID)
	}

	closestNode := contactList[0]

	for {
		unprobedContacts := kademlia.GetNextUnprobedContacts(contactList)
		if len(unprobedContacts) == 0 {
			return GetAllContactsFromContactList(contactList), Contact{}, nil
		}

		var contactFoundDataOn Contact
		var foundData []byte

		contactList, contactFoundDataOn, foundData = kademlia.SendFindNodeMessagesToUnprobedContacts(contactList, target, hash, unprobedContacts)
		fmt.Println("DEBUG: ContactList after sending messages", contactList)

		if foundData != nil {
			fmt.Println("Node lookup complete, data found")
			return GetAllContactsFromContactList(contactList), contactFoundDataOn, foundData
		}

		newClosestNode := contactList[0]

		if closestNode.Contact.ID.Equals(newClosestNode.Contact.ID) {
			unprobedKClosest := kademlia.GetNextUnprobedContacts(contactList)
			if CountProbedInContactList(contactList) >= k || len(unprobedKClosest) == 0 {
				break
			} else {
				nextContactList, _, _ := kademlia.SendFindNodeMessagesToUnprobedContacts(contactList, target, hash, unprobedKClosest)
				contactList = nextContactList
			}
		} else {
			closestNode = newClosestNode
		}
	}

	fmt.Println("Node lookup complete")
	return GetAllContactsFromContactList(contactList), Contact{}, nil
}

func AddContactToContactList(contactList []ContactDistance, newContact Contact, targetID *KademliaID) []ContactDistance {
	for _, item := range contactList {
		if item.Contact.ID.Equals(newContact.ID) {
			return contactList
		}
	}

	newDistance := newContact.ID.CalcDistance(targetID)
	newItem := ContactDistance{
		Contact:          newContact,
		DistanceToTarget: newDistance,
		Probed:           false,
	}

	contactList = append(contactList, newItem)

	sort.Slice(contactList, func(i, j int) bool {
		return contactList[i].DistanceToTarget.Less(contactList[j].DistanceToTarget)
	})

	if len(contactList) > k && k > 0 {
		contactList = contactList[:k]
	}

	return contactList
}

func (kademlia *Kademlia) GetNextUnprobedContacts(contactList []ContactDistance) []ContactDistance {
	unprobedContacts := make([]ContactDistance, 0, alpha)

	for _, item := range contactList {
		if !item.Probed && !item.Contact.ID.Equals(kademlia.RoutingTable.me.ID) {
			unprobedContacts = append(unprobedContacts, item)
		}
		if len(unprobedContacts) == alpha {
			break
		}
	}

	return unprobedContacts
}

func GetAllContactsFromContactList(contactList []ContactDistance) []Contact {
	contacts := make([]Contact, len(contactList))
	for i, item := range contactList {
		contacts[i] = item.Contact
	}
	return contacts
}

func (kademlia *Kademlia) SendFindNodeMessagesToUnprobedContacts(contactList []ContactDistance, target *Contact, hash string, unprobedContacts []ContactDistance) ([]ContactDistance, Contact, []byte) {
	contactResponses := make(chan Contact, alpha)
	dataResponses := make(chan []byte, alpha)
	contactWithData := make(chan Contact, alpha)

	go kademlia.sendFindNodeQueries(unprobedContacts, target, hash, contactResponses, dataResponses, contactWithData)

	defer close(contactResponses)
	defer close(dataResponses)
	defer close(contactWithData)

	fmt.Println("DEBUG: Completed contact probing")

	foundContact, foundData := handleFoundData(dataResponses, contactWithData)
	if foundData != nil {
		return contactList, foundContact, foundData
	}

	updatedContactList := kademlia.updateContactListWithReceivedContacts(contactList, target, contactResponses)
	updatedContactList = markContactsAsProbed(updatedContactList, unprobedContacts)

	return updatedContactList, Contact{}, nil
}

func (kademlia *Kademlia) updateContactListWithReceivedContacts(contactList []ContactDistance, target *Contact, contactsChan chan Contact) []ContactDistance {
	for contact := range contactsChan {
		found := false
		for i := range contactList {
			if contactList[i].Contact.ID.Equals(contact.ID) {
				contactList[i].Contact = contact
				found = true
				break
			}
		}

		if !found {
			contactList = AddContactToContactList(contactList, contact, target.ID)
		}
	}
	return contactList
}

func markContactsAsProbed(contactList []ContactDistance, unprobedContacts []ContactDistance) []ContactDistance {
	unprobedMap := make(map[string]bool)
	for _, contact := range unprobedContacts {
		unprobedMap[contact.Contact.ID.String()] = true
	}

	for i, contactItem := range contactList {
		if unprobedMap[contactItem.Contact.ID.String()] {
			contactList[i].Probed = true
		}
	}

	return contactList
}

func CountProbedInContactList(contactList []ContactDistance) int {
	probedCount := 0
	for _, contact := range contactList {
		if contact.Probed {
			probedCount++
		}
	}
	return probedCount
}

func (kademlia *Kademlia) sendFindNodeQueries(unprobedContacts []ContactDistance, target *Contact, hash string, contactResponses chan Contact, dataResponses chan []byte, contactWithDataChan chan Contact) {
	var wg sync.WaitGroup

	queryType := "FIND_NODE"
	if hash != "" {
		queryType = "FIND_DATA"
	}

	queryFunc := func(contact Contact) {
		defer wg.Done()
		if queryType == "FIND_NODE" {
			kademlia.Network.SendFindContactMessage(contact, target, contactResponses)
		} else {
			kademlia.Network.SendFindDataMessage(contact, hash, dataResponses, contactWithDataChan)
		}
	}

	for _, contactItem := range unprobedContacts {
		wg.Add(1)
		go queryFunc(contactItem.Contact)
	}

	wg.Wait()
}

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

	if _, found := (*kademlia.Data)[key]; found {
		fmt.Println("Data could not be modified, already exist data with this key.")
		return
	}

	(*kademlia.Data)[key] = data
	fmt.Println(key)

	storeID := NewKademliaID(key)
	contact := NewContact(storeID, "")

	contacts, err := kademlia.LookupContact(&contact)
	if err != nil {
		fmt.Printf("Failed to lookup contacts for storage: %v\n", err)
		return
	}

	fmt.Println("Stored data with key:", key)

	for _, contact := range contacts {
		fmt.Printf("Stored data at contact: %s with key: %s\n", contact.Address, key)
		go kademlia.Network.SendStoreMessage(data, &contact)
	}
}

func (kademlia *Kademlia) Get(hash string) ([]byte, []Contact, error) {
	data, closestContacts, err := kademlia.LookupData(hash)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to lookup data: %w", err)
	}

	if data != nil {
		fmt.Println("Data found locally.")
		return data, nil, nil
	}

	fmt.Println("Data not found locally. Closest contacts returned.")
	return nil, closestContacts, nil
}
