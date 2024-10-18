package kademlia

import (
	"fmt"
	"testing"
	"time"
)

func TestJoinNetwork(t *testing.T) {
	// Setup real network data
	me := NewContact(NewKademliaID("c3075781c1b20059c02c7a0bdea2aabfa9e9537022e0b520eac716979264c8aa"), "localhost:3002") //We need to define a node that are to join the network.
	network := Network{
		&me,                 //We pass the pointer to me
		NewRoutingTable(me), //Not a defined as a pointer here due to NewRoutingTable returning a pointer object of RoutingTable
	}
	go network.Listen("0.0.0.0", 3002)

	me2 := NewContact(NewKademliaID("93a800ff24d53aaa51e934175ae01ae5dd01e97b55dc102c3a79f5a0d0a2f9b8"), "localhost:3002") //We need to define a node that are to join the network.

	network.SendJoinMessage(&me)
	network.SendJoinMessage(&me2)

	fmt.Printf("Node Address111: %s\n", me.Address)
	fmt.Printf("Node Address222: %s\n", me2.Address)

	// Get the closest contacts to verify the join process
	closestContacts := network.RoutingTable.FindClosestContacts(NewKademliaID("c3075781c1b20059c02c7a0bdea2aabfa9e9537022e0b520eac716979264c8aa"), 1)

	// Check if the node successfully joined the network
	if closestContacts[0].ID.String() != me.ID.String() {
		t.Errorf("Failed to join network. Expected closest contact to be %s, got %s", me.ID.String(), closestContacts[0].ID.String())
	}
}
func TestSendPingMessage(t *testing.T) {
	me := NewContact(NewKademliaID("d4838ebed2c547b6ab87e1f70b789d4f94ce7a85622c47143ade0d3a7ce4d0e4"), "localhost:3001")
	network := &Network{
		Self:         &me,
		RoutingTable: NewRoutingTable(me),
	}

	go network.Listen("0.0.0.0", 3001)
	time.Sleep(1 * time.Second)

	// Send Ping Message and wait for response
	network.SendPingMessage(&me)

	// Sleep to allow time for handling response
	time.Sleep(1 * time.Second)
}

func TestSendStoreMessage(t *testing.T) {
	me := NewContact(NewKademliaID("ffffffff00000000000000000000000000000000"), "localhost:8000")
	target := NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8001")

	network := &Network{
		Self:         &me,
		RoutingTable: NewRoutingTable(me),
	}

	go network.Listen("0.0.0.0", 8000)
	time.Sleep(1 * time.Second)

	network2 := &Network{
		Self:         &target,
		RoutingTable: NewRoutingTable(target),
	}

	go network2.Listen("0.0.0.0", 8001)
	time.Sleep(1 * time.Second)

	// Send store message to target
	data := []byte("hello world")
	network.SendStoreMessage(data, &target)

	// Wait for the store message to be processed
	time.Sleep(2 * time.Second)
}
