package kademlia

import "testing"

func TestSendPingMessage(t *testing.T) {
	// Setup real network data
	me      := NewContact(NewKademliaID("41db94eb909853000c950012145f9c6479d63a9a97c002803ff93be4b38a9211"), "localhost:3001") //We need to define a node that are to join the network.
    network := Network{
		&me, 							   //We pass the pointer to me
		NewRoutingTable(me), //Not a defined as a pointer here due to NewROutingTable returning a pointer object of RoutingTable
    }

	// Provide a 40-character hex string (which is 20 bytes) for the KademliaID
	contact := &Contact{
		ID:      NewKademliaID("41db94eb909853000c950012145f9c6479d63a9a97c002803ff93be4b38a9211"), // 40-character hex string it will cut the id from the back to exactly 40 chars.
		Address: "127.0.0.1",                                                 // Use Address instead of IP
	}

	// Call the real SendPingMessage function and check for proper behavior
	network.SendPingMessage(contact)
}




// 
// 
// 

// SLUTADE ATT KOLLA OM MAN KUNDE FÅ TVÅ NODER I SAMMA NÄTVERK OCH SE I DOCKER DESKTOP OM DESSA LÅG I SAMMA NÄTVERK ELLER VIA TERMINALEN KOLLA 
// DETTA, MEN JAG STÖTTE PÅ PROBLEM DÄR JAG BUILDAR OCH STARTAR UPP DOCKER VIA TERMINALEN ATT ALLA NODER STARTAR OCH SEDAN EXITAR EFTER 10 SEK OCH 
// DÅ KAN MAN JU INTE ANVÄNDA DEM FÖR DE INTE ÄR PÅ.

// func TestNodeLookup(t *testing.T){

// }

// //TestJoinNetwork tests the JoinNetwork function
// func TestJoinNetwork(t *testing.T) {
//     // Setup real network data
// 	me      := NewContact(NewKademliaID("41db94eb909853000c950012145f9c6479d63a9a97c002803ff93be4b38a9211"), "localhost:3005") //We need to define a node that are to join the network.
//     network := Network{
// 		&me, 				 //We pass the pointer to me
// 		NewRoutingTable(me), //Not a defined as a pointer here due to NewRoutingTable returning a pointer object of RoutingTable
//     }

// 	me2      := NewContact(NewKademliaID("75111427089cba531b716e25c0ffefece4f7832e65d61248bde799c0ac6f8b3a"), "localhost:3005") //We need to define a node that are to join the network.

// 	err  := network.JoinNetwork(&me)
// 	err2 := network.JoinNetwork(&me2)
	
// 	if err != nil {
// 		t.Errorf("Failed to join network: %v", err)
// 		return
// 	}
// 	if err2 != nil {
// 		t.Errorf("Failed to join network: %v", err)
// 		return
// 	}
	
// 	// Get the closest contacts to verify the join process
// 	closestContacts := network.RoutingTable.FindClosestContacts(NewKademliaID("41db94eb909853000c950012145f9c6479d63a9a97c002803ff93be4b38a9211"), 1)

// 	// Check if the node successfully joined the network
// 	if closestContacts[0].ID.String() != me.ID.String() {
// 		t.Errorf("Failed to join network. Expected closest contact to be %s, got %s", me.ID.String(), closestContacts[0].ID.String())
// 	}
// }