package kademlia

import (
	"fmt"
	"testing"
)

func TestSendPingMessage(t *testing.T) {
	// Setup real network data
	me := NewContact(NewKademliaID("d4838ebed2c547b6ab87e1f70b789d4f94ce7a85622c47143ade0d3a7ce4d0e4"), "localhost:3001") //We need to define a node that are to join the network.
	network := Network{
		&me,                 //We pass the pointer to me
		NewRoutingTable(me), //Not a defined as a pointer here due to NewROutingTable returning a pointer object of RoutingTable
	}

	go network.Listen("0.0.0.0", 3001)
	// Call the real SendPingMessage function and check for proper behavior
	network.SendPingMessage(&me)
}

func TestNodeLookup(t *testing.T) {

}

// TestJoinNetwork tests the JoinNetwork function
func TestJoinNetwork(t *testing.T) {
	// Setup real network data
	me := NewContact(NewKademliaID("eacd8eaf02d284c497494f7473bf49991b250bac94fae100acd7dfee00fd6fe5"), "localhost:3002") //We need to define a node that are to join the network.
	network := Network{
		&me,                 //We pass the pointer to me
		NewRoutingTable(me), //Not a defined as a pointer here due to NewRoutingTable returning a pointer object of RoutingTable
	}
	go network.Listen("0.0.0.0", 3002)

	me2 := NewContact(NewKademliaID("9f3e804b394c86bd95b8ae1d77a1ec4ba5ce5de51faa3b7d2681419b2089f3d6"), "localhost:3002") //We need to define a node that are to join the network.

	network.JoinNetwork(&me)
	network.JoinNetwork(&me2)

	fmt.Printf("Node Address111: %s\n", me.Address)
	fmt.Printf("Node Address222: %s\n", me2.Address)

	// Get the closest contacts to verify the join process
	closestContacts := network.RoutingTable.FindClosestContacts(NewKademliaID("eacd8eaf02d284c497494f7473bf49991b250bac94fae100acd7dfee00fd6fe5"), 1)

	// Check if the node successfully joined the network
	if closestContacts[0].ID.String() != me.ID.String() {
		t.Errorf("Failed to join network. Expected closest contact to be %s, got %s", me.ID.String(), closestContacts[0].ID.String())
	}
}
