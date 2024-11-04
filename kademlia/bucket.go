package kademlia

import (
	"container/list"
	"fmt"
)

// bucket definition
// contains a List
type bucket struct {
	list *list.List
}

// newBucket returns a new instance of a bucket
func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	return bucket
}

// AddContact adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *bucket) AddContact(contact Contact) (bool, *Contact) {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	//if contact not in bucket
	if element == nil {
		//add to front i there is space
		if bucket.list.Len() < bucketSize {
			bucket.list.PushFront(contact)
			return false, nil
		}

		//full bucket, return last contact
		lastContact := bucket.list.Back().Value.(Contact)
		return true, &lastContact

	} else {
		//moving found contact to front of list
		bucket.list.MoveToFront(element)
		return false, nil
	}
}

// GetContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

// Len return the size of the bucket
func (bucket *bucket) Len() int {
	return bucket.list.Len()
}

// RemoveContact deletes existing contact from bucket
func (bucket *bucket) RemoveContact(contact *Contact) {
	current := bucket.list.Front()

	for current != nil {
		existingContact := current.Value.(Contact)

		if existingContact.ID.Equals(contact.ID) {
			bucket.list.Remove(current)
			return
		}

		current = current.Next()
	}
}

func (bucket *bucket) PrintIPs() {
	for element := bucket.list.Front(); element != nil; element = element.Next() {
		contact := element.Value.(Contact)
		fmt.Println("Address: " + contact.Address + " ID: " + contact.ID.String())
	}
}
