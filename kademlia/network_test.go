package kademlia

import (
	"testing"
	"fmt"
)

func TestSendPingMessage(t *testing.T) {
	// Setup real network data
	me      := NewContact(NewKademliaID("ffa5ffe6b55a05ccf2b8e3f601d199e178ec5538e67b52d73e28a324039c5a58"), "localhost:3001") //We need to define a node that are to join the network.
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
	me      := NewContact(NewKademliaID("d5023b0433620a6a1a38715f14601dd3c1553d763fc8bee9310c7ec6fcf8d6a3"), "localhost:3002") //We need to define a node that are to join the network.
    network := Network{
		&me, 				 //We pass the pointer to me
		NewRoutingTable(me), //Not a defined as a pointer here due to NewRoutingTable returning a pointer object of RoutingTable
    }
	go network.Listen("0.0.0.0", 3002)

	me2      := NewContact(NewKademliaID("ffa5ffe6b55a05ccf2b8e3f601d199e178ec5538e67b52d73e28a324039c5a58"), "localhost:3002") //We need to define a node that are to join the network.

	network.JoinNetwork(&me)
	network.JoinNetwork(&me2)
	


    fmt.Printf("Node Address111: %s\n", me.Address)
    fmt.Printf("Node Address222: %s\n", me2.Address)




	// Get the closest contacts to verify the join process
	closestContacts := network.RoutingTable.FindClosestContacts(NewKademliaID("d5023b0433620a6a1a38715f14601dd3c1553d763fc8bee9310c7ec6fcf8d6a3"), 1)



	
	// Check if the node successfully joined the network
	if closestContacts[0].ID.String() != me.ID.String() {
		t.Errorf("Failed to join network. Expected closest contact to be %s, got %s", me.ID.String(), closestContacts[0].ID.String())
	}
}