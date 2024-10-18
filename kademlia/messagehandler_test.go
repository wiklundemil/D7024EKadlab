package kademlia

import (
	"encoding/json"
	"testing"
)

func TestHandlePingMessage(t *testing.T) {
	me := NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000")
	network := &Network{
		Self:         &me,
		RoutingTable: NewRoutingTable(me),
	}

	// Create a PING message
	ping := Message{
		MessageType: "PING",
		Content:     "Ping from localhost:8001",
		Sender:      me,
	}

	byteStream, _ := json.Marshal(ping)
	response, err := network.HandleMessage(byteStream)
	if err != nil {
		t.Fatalf("Failed to handle message: %v", err)
	}

	var msg Message
	json.Unmarshal(response, &msg)

	if msg.MessageType != "PING_ACK" {
		t.Fatalf("Expected PING_ACK but got %s", msg.MessageType)
	}
}

func TestHandleJoinMessage(t *testing.T) {
	me := NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000")
	network := &Network{
		Self:         &me,
		RoutingTable: NewRoutingTable(me),
	}

	// Create a JOIN message
	join := Message{
		MessageType: "JOIN",
		Content:     "localhost:8001",
		Sender:      me,
	}

	byteStream, _ := json.Marshal(join)
	response, err := network.HandleMessage(byteStream)
	if err != nil {
		t.Fatalf("Failed to handle message: %v", err)
	}

	var msg Message
	json.Unmarshal(response, &msg)

	if msg.MessageType != "JOIN_ACK" {
		t.Fatalf("Expected JOIN_ACK but got %s", msg.MessageType)
	}
}

func TestHandleStoreMessage(t *testing.T) {
	me := NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000")
	network := &Network{
		Self:         &me,
		RoutingTable: NewRoutingTable(me),
	}

	// Create a STORE message
	store := Message{
		MessageType: "STORE",
		Content:     "Sample data",
		Sender:      me,
	}

	byteStream, _ := json.Marshal(store)
	response, err := network.HandleMessage(byteStream)
	if err != nil {
		t.Fatalf("Failed to handle message: %v", err)
	}

	var msg Message
	json.Unmarshal(response, &msg)

	if msg.MessageType != "STORE_ACK" {
		t.Fatalf("Expected STORE_ACK but got %s", msg.MessageType)
	}
}
