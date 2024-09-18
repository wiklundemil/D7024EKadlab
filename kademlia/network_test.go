package kademlia

import "testing"

func TestSendPingMessage(t *testing.T) {
	// Setup real network data
	me      := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3001") //We need to define a node that are to join the network.
    network := Network{
		&me, 							   //We pass the pointer to me
		NewRoutingTable(me), //Not a defined as a pointer here due to NewROutingTable returning a pointer object of RoutingTable
    }

	// Provide a 40-character hex string (which is 20 bytes) for the KademliaID
	contact := &Contact{
		ID:      NewKademliaID("7461726765742d6e6f64652d69646d6f636b69646"), // 40-character hex string it will cut the id from the back to exactly 40 chars.
		Address: "127.0.0.1",                                                 // Use Address instead of IP
	}

	// Call the real SendPingMessage function and check for proper behavior
	network.SendPingMessage(contact)
}


//TestJoinNetwork tests the JoinNetwork function
func TestJoinNetwork(t *testing.T) {
    // Setup real network data
	me      := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:3005") //We need to define a node that are to join the network.
    network := Network{
		&me, 							   //We pass the pointer to me
		NewRoutingTable(me), //Not a defined as a pointer here due to NewROutingTable returning a pointer object of RoutingTable
    }
	err := network.JoinNetwork(&me)
	
	if err != nil {
		t.Errorf("Failed to join network: %v", err)
		return
	}
	
	// Get the closest contacts to verify the join process
	closestContacts := network.RoutingTable.FindClosestContacts(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), 1)

	// Check if the node successfully joined the network
	if closestContacts[0].ID.String() != me.ID.String() {
		t.Errorf("Failed to join network. Expected closest contact to be %s, got %s", me.ID.String(), closestContacts[0].ID.String())
	}
}