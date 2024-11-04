package kademlia

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

const bucketTestSize = 5

func TestNewBucket(t *testing.T) {
	bucket := newBucket()

	// Check if the bucket itself is nil
	if bucket == nil {
		t.Error("Failed to create a new bucket")
		return
	}

	// Check if the bucket list is initialized
	if bucket.list == nil {
		t.Error("Bucket list initialization failed")
		return
	}
}

func TestAddContact(t *testing.T) {
	bucket := newBucket()
	contact := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")

	// Adding a new contact to empty bucket
	bucketIsFull, lastContact := bucket.AddContact(contact)
	if bucketIsFull {
		t.Error("Expected bucket to not be full after adding the first contact")
	}
	if lastContact != nil {
		t.Error("Expected lastContact to be nil when the bucket is not full")
	}

	// Test moving an existing contact to the front
	bucket.AddContact(contact)
	if bucket.list.Front().Value.(Contact).ID != contact.ID {
		t.Error("Expected existing contact to be moved to the front of the list")
	}

	// Fill the bucket with new contacts
	for i := 0; i < bucketTestSize; i++ {
		bucket.AddContact(NewContact(NewRandomKademliaID(), "127.0.0.1:8000"))
	}

	// Add one last contact to fill up the bucket
	newContact := NewContact(NewRandomKademliaID(), "127.0.0.1:9000")
	bucketIsFull, lastContact = bucket.AddContact(newContact)
	if !bucketIsFull {
		t.Error("Expected bucket to be full after adding a contact beyond the capacity")
	}
	if lastContact == nil {
		t.Error("Expected lastContact to be the last item in the bucket when full")
	}
}

func TestGetContactAndCalcDistance(t *testing.T) {
	bucket := newBucket()
	targetID := NewRandomKademliaID()
	contact := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")
	bucket.AddContact(contact)
	contacts := bucket.GetContactAndCalcDistance(targetID)

	// Verify that just one contact is returned
	if len(contacts) != 1 {
		t.Error("Expected 1 contact in the bucket")
	}

	// Check that returned contact matches added contact
	if contacts[0].ID != contact.ID {
		t.Error("Returned contact does not match the added contact")
	}
}

func TestLen(t *testing.T) {
	bucket := newBucket()
	if length := bucket.Len(); length != 0 {
		t.Error("Expected bucket to be of length 0")
	}

	// Add contact to bucket
	contact := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")
	bucket.AddContact(contact)

	// Verify that the bucket length is exactly 1
	if length := bucket.Len(); length != 1 {
		t.Error("Expected bucket to be of length 1 after adding a contact")
	}
}

func TestRemoveContact(t *testing.T) {
	bucket := newBucket()
	contact := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")
	bucket.AddContact(contact)

	if bucket.list.Len() == 0 {
		t.Error("Failed to add contact to bucket")
	}

	// Remove contact
	bucket.RemoveContact(&contact)

	if bucket.list.Len() != 0 {
		t.Error("Contact was not removed from bucket as expected")
	}
}

func TestPrintIPsBucket(t *testing.T) {
	// Set up a new bucket and add contacts
	bucket := newBucket()
	contact1 := NewContact(NewRandomKademliaID(), "127.0.0.1:8000")
	contact2 := NewContact(NewRandomKademliaID(), "127.0.0.1:8001")
	bucket.AddContact(contact1)
	bucket.AddContact(contact2)

	// Redirect stdout to capture output
	r, w, _ := os.Pipe()
	defer r.Close()
	oldStdout := os.Stdout
	os.Stdout = w

	bucket.PrintIPs()

	// Restore stdout and close the writer
	w.Close()
	os.Stdout = oldStdout

	// Capture and read the output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify the output contains both contacts' address and ID
	for _, contact := range []Contact{contact1, contact2} {
		if !strings.Contains(output, contact.Address) || !strings.Contains(output, contact.ID.String()) {
			t.Errorf("Expected output to contain contact's address (%s) and ID (%s), got: %s", contact.Address, contact.ID.String(), output)
		}
	}
}
