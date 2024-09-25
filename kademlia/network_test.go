package kademlia

import (
	"testing"
	"fmt"
)

func TestSendPingMessage(t *testing.T) {
	// Setup real network data
	me      := NewContact(NewKademliaID("1f57442a7b6b94b1a0e175b63a727272de8940de05d20e26c4e19db0044b92b9"), "localhost:3001") //We need to define a node that are to join the network.
    network := Network{
		&me, 							   //We pass the pointer to me
		NewRoutingTable(me), //Not a defined as a pointer here due to NewROutingTable returning a pointer object of RoutingTable
    }

	go network.Listen("0.0.0.0", 3001)
	// Call the real SendPingMessage function and check for proper behavior
	network.SendPingMessage(&me)
}

func TestNodeLookup(t *testing.T){

}

//TestJoinNetwork tests the JoinNetwork function
func TestJoinNetwork(t *testing.T) {
    // Setup real network data
	me      := NewContact(NewKademliaID("9ee16f7ca8e6490efd9f6c7f8934becd4c96899572ad40e207d8074e68de5eb7"), "localhost:3002") //We need to define a node that are to join the network.
    network := Network{
		&me, 				 //We pass the pointer to me
		NewRoutingTable(me), //Not a defined as a pointer here due to NewRoutingTable returning a pointer object of RoutingTable
    }
	go network.Listen("0.0.0.0", 3002)

	me2      := NewContact(NewKademliaID("6ff6fb44432bf5f1133f056ec131fe9bb7caa29aa84f48da6a90b3804d5cba36"), "localhost:3002") //We need to define a node that are to join the network.

	network.JoinNetwork(&me)
	network.JoinNetwork(&me2)

    fmt.Printf("Node Address111: %s\n", me.Address)
    fmt.Printf("Node Address222: %s\n", me2.Address)

	// Get the closest contacts to verify the join process
	closestContacts := network.RoutingTable.FindClosestContacts(NewKademliaID("9ee16f7ca8e6490efd9f6c7f8934becd4c96899572ad40e207d8074e68de5eb7"), 1)

	// Check if the node successfully joined the network
	if closestContacts[0].ID.String() != me.ID.String() {
		t.Errorf("Failed to join network. Expected closest contact to be %s, got %s", me.ID.String(), closestContacts[0].ID.String())
	}
}