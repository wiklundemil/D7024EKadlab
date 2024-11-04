package kademlia

import (
	"fmt"
	"net"
	"sort"
	"sync"
)

type Kademlia struct {
	RoutingTable  *RoutingTable
	Network       *Network
	Data          *map[string][]byte
	ActionChannel chan Action
}

type Action struct {
	Action   string
	Target   *Contact
	Hash     string
	Data     []byte
	SenderId *KademliaID
	SenderIp string
}

type ContactListItem struct {
	Contact          Contact
	DistanceToTarget *KademliaID
	Probed           bool
}

const (
	alpha = 3 
	k     = 5 
)

func NewKademlia(rTable *RoutingTable, conn net.PacketConn) *Kademlia {
	netLayer := NewNetwork(conn)
	store := make(map[string][]byte)
	actionPipe := make(chan Action)
	return &Kademlia{RoutingTable: rTable, Network: netLayer, Data: &store, ActionChannel: actionPipe}
}
func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	closestContacts := kademlia.RoutingTable.FindClosestContacts(target.ID, k)
	return closestContacts
}

func (kademlia *Kademlia) LookupData(hash string) ([]byte, []Contact) {
	if value, found := (*kademlia.Data)[hash]; found {
		return value, nil
	}

	targetContact := NewContact(NewKademliaID(hash), "")
	nearestContacts := kademlia.LookupContact(&targetContact)
	return nil, nearestContacts
}

func (kademlia *Kademlia) Store(hash string, data []byte) {
	(*kademlia.Data)[hash] = data
}

func (kademlia *Kademlia) NodeLookup(target *Contact, hash string) ([]Contact, Contact, []byte) {
	initialContacts := kademlia.RoutingTable.FindClosestContacts(target.ID, alpha)
	var candidateList []ContactListItem
	for _, contact := range initialContacts {
		candidateList = UpdateContactList(candidateList, contact, target.ID)
	}

	nearestContact := candidateList[0]

	for {
		remainingUnprobed := kademlia.GetAlpha(candidateList)
		if len(remainingUnprobed) == 0 {
			return GetAllContactsFromContactList(candidateList), Contact{}, nil
		}

		unprobedNodes := kademlia.GetAlpha(candidateList)
		var dataProvider Contact
		var retrievedData []byte

		candidateList, dataProvider, retrievedData = kademlia.SendAlphaFindNodeMessages(candidateList, target, hash, unprobedNodes)

		if retrievedData != nil {
			fmt.Println("Node lookup complete: data found")
			return GetAllContactsFromContactList(candidateList), dataProvider, retrievedData
		}

		newNearestContact := candidateList[0]

		if nearestContact.Contact.ID.Equals(newNearestContact.Contact.ID) {
			moreUnprobed := kademlia.GetAlpha(candidateList)
			if CountProbedInContactList(candidateList) >= k || len(moreUnprobed) == 0 {
				break
			} else {
				closestUnprobed := kademlia.GetAlphaFromKClosest(candidateList, target)
				updatedList, _, _ := kademlia.SendAlphaFindNodeMessages(candidateList, target, hash, closestUnprobed)
				candidateList = updatedList
			}
		} else {
			nearestContact = newNearestContact
		}
	}
	fmt.Println("Node lookup completed without finding data")
	return GetAllContactsFromContactList(candidateList), Contact{}, nil
}

func (kademlia *Kademlia) UpdateRT(id *KademliaID, ip string) {
	newContact := NewContact(id, ip)
	if !newContact.ID.Equals(kademlia.RoutingTable.Me.ID) {
		fmt.Printf("Inserting contact to routing table with ID: %s and IP: %s on %s\n", newContact.ID.String(), newContact.Address, kademlia.RoutingTable.Me.Address)
		newContact.CalcDistance(kademlia.RoutingTable.Me.ID)

		isBucketFull, previousContact := kademlia.RoutingTable.AddContact(newContact)
		if isBucketFull {
			if kademlia.Network.SendPingMessage(&kademlia.RoutingTable.Me, previousContact) {
				fmt.Println("Previous contact is responsive, discarding the new contact")
			} else {
				fmt.Println("Previous contact is unresponsive, replacing with the new contact")
				kademlia.RoutingTable.RemoveContact(previousContact)
				kademlia.RoutingTable.AddContact(newContact)
			}
		}
	}
}

func UpdateContactList(contactList []ContactListItem, newContact Contact, target *KademliaID) []ContactListItem {
	for _, entry := range contactList {
		if entry.Contact.ID.Equals(newContact.ID) {
			return contactList
		}
	}

	distanceToTarget := newContact.ID.CalcDistance(target)
	entry := ContactListItem{Contact: newContact, DistanceToTarget: distanceToTarget, Probed: false}
	contactList = append(contactList, entry)

	sort.Slice(contactList, func(i, j int) bool {
		return contactList[i].DistanceToTarget.Less(contactList[j].DistanceToTarget)
	})

	if len(contactList) > k {
		contactList = contactList[:k]
	}
	return contactList
}

func GetAllContactsFromContactList(contactList []ContactListItem) []Contact {
	contactsList := make([]Contact, 0, len(contactList))
	for _, entry := range contactList {
		contactsList = append(contactsList, entry.Contact)
	}
	return contactsList
}

func (kademlia *Kademlia) probeContacts(unprobedContacts []ContactListItem, target *Contact, hashKey string, contactChannel chan Contact, dataChannel chan []byte, contactDataChannel chan Contact) {
	var waitGroup sync.WaitGroup

	for _, contactItem := range unprobedContacts {
		waitGroup.Add(1)
		go func(contactInfo Contact) {
			defer waitGroup.Done()
			if hashKey == "" {
				kademlia.findContact(contactInfo, target, contactChannel, dataChannel, contactDataChannel)
			} else {
				kademlia.findData(contactInfo, hashKey, dataChannel, contactDataChannel)
			}
		}(contactItem.Contact)
	}
	waitGroup.Wait()
}

func closeChannels(contactChannel chan Contact, dataChannel chan []byte, foundContactChannel chan Contact) {
	if contactChannel != nil {
		close(contactChannel)
	}
	if dataChannel != nil {
		close(dataChannel)
	}
	if foundContactChannel != nil {
		close(foundContactChannel)
	}
}

func handleFoundData(dataChannel chan []byte, foundContactChannel chan Contact) (Contact, []byte) {
	select {
	case retrievedData := <-dataChannel:
		if retrievedData != nil {
			associatedContact := <-foundContactChannel
			if associatedContact.ID != nil {
				return associatedContact, retrievedData
			}
		}
	default:
	}
	return Contact{}, nil
}

func (kademlia *Kademlia) updateContactListWithContacts(contactList []ContactListItem, target *Contact, contactStream chan Contact) []ContactListItem {
	for receivedContact := range contactStream {
		isUpdated := false
		for index, item := range contactList {
			if item.Contact.ID.Equals(receivedContact.ID) {
				contactList[index].Contact = receivedContact
				isUpdated = true
				break
			}
		}
		if !isUpdated {
			contactList = UpdateContactList(contactList, receivedContact, target.ID)
		}
	}
	return contactList
}

func markProbedContacts(contactList []ContactListItem, unmarkedContacts []ContactListItem) []ContactListItem {
	for idx, ContactlistItem := range contactList {
		for _, unmarkedItem := range unmarkedContacts {
			if ContactlistItem.Contact.ID.Equals(unmarkedItem.Contact.ID) {
				contactList[idx].Probed = true
				break
			}
		}
	}
	return contactList
}
func (kademlia *Kademlia) SendAlphaFindNodeMessages(contactList []ContactListItem, target *Contact, hash string, unqueriedNodes []ContactListItem) ([]ContactListItem, Contact, []byte) {
	nodeChannel := make(chan Contact, alpha*k)
	dataChannel := make(chan []byte, alpha*k)
	foundContactChannel := make(chan Contact, alpha*k)

	kademlia.probeContacts(unqueriedNodes, target, hash, nodeChannel, dataChannel, foundContactChannel)

	closeChannels(nodeChannel, dataChannel, foundContactChannel)

	discoveredContact, foundData := handleFoundData(dataChannel, foundContactChannel)
	if foundData != nil {
		return contactList, discoveredContact, foundData
	}

	contactList = kademlia.updateContactListWithContacts(contactList, target, nodeChannel)
	contactList = markProbedContacts(contactList, unqueriedNodes)

	return contactList, Contact{}, nil
}

func (kademlia *Kademlia) findContact(contact Contact, target *Contact, nodeChannel chan Contact, responseDataChan chan []byte, responseContactChan chan Contact) {
	retrievedContacts, err := kademlia.Network.SendFindContactMessage(&kademlia.RoutingTable.Me, &contact, target)
	if err != nil {
		fmt.Printf("Error occurred while sending FIND_NODE message: %v\n", err)
		return
	}

	for _, retrievedContact := range retrievedContacts {
		select {
		case nodeChannel <- retrievedContact:
			responseDataChan <- nil
			responseContactChan <- Contact{}
		default:
			fmt.Printf("Channel buffer full, could not send contact: %s\n", retrievedContact.String())
		}
	}
}

func (kademlia *Kademlia) findData(contact Contact, hashValue string, dataChannel chan []byte, responseContactChan chan Contact) {
	_, retrievedData, err := kademlia.Network.SendFindDataMessage(&kademlia.RoutingTable.Me, &contact, hashValue)
	if err != nil {
		fmt.Printf("Error during FIND_DATA message: %v\n", err)
		return
	}

	if retrievedData != nil {
		dataChannel <- retrievedData
		responseContactChan <- contact
	}
}

func (kademlia *Kademlia) GetAlpha(contactList []ContactListItem) []ContactListItem {
	var unprobedContacts []ContactListItem

	for _, contactItem := range contactList {
		if !contactItem.Probed && !contactItem.Contact.ID.Equals(kademlia.RoutingTable.Me.ID) {
			unprobedContacts = append(unprobedContacts, contactItem)
		}
	}

	if len(unprobedContacts) <= alpha {
		return unprobedContacts
	}
	return unprobedContacts[:alpha]
}
func (kademlia *Kademlia) GetAlphaFromKClosest(candidateList []ContactListItem, destination *Contact) []ContactListItem {
	var untestedNodes []ContactListItem
	closestContacts := kademlia.RoutingTable.FindClosestContacts(destination.ID, k)
	for _, contact := range closestContacts {
		for _, candidate := range candidateList {
			if contact.ID.Equals(candidate.Contact.ID) || len(untestedNodes) >= alpha {
				continue
			} else {
				untestedNodes = append(untestedNodes, ContactListItem{contact, contact.ID.CalcDistance(destination.ID), false})
			}
		}
	}

	if len(untestedNodes) < alpha {
		return untestedNodes
	}
	return untestedNodes[:alpha]
}

func (kademlia *Kademlia) ListenActionChannel() {
	for {
		currentAction := <-kademlia.ActionChannel
		switch currentAction.Action {
		case "UpdateRT":
			kademlia.UpdateRT(currentAction.SenderId, currentAction.SenderIp)
		case "Store":
			kademlia.Store(currentAction.Hash, currentAction.Data)
		case "LookupContact":
			closestNodes := kademlia.LookupContact(currentAction.Target)
			lookupResponse := Response{
				ClosestContacts: closestNodes,
			}
			kademlia.Network.responseChan <- lookupResponse
		case "LookupData":
			foundData, nodesList := kademlia.LookupData(currentAction.Hash)
			dataResponse := Response{
				Data:            foundData,
				ClosestContacts: nodesList,
			}
			kademlia.Network.responseChan <- dataResponse
		case "PRINT":
			kademlia.RoutingTable.PrintIPs()
		}
	}
}

func CountProbedInContactList(contactList []ContactListItem) int {
	probedCount := 0
	for _, contact := range contactList {
		if contact.Probed {
			probedCount += 1
		}
	}
	return probedCount
}
