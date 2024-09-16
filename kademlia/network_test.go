package kademlia

import "testing"

func TestSendPingMessage(t *testing.T) {
	// Setup real network data
	network := &Network{
		SelfID: "self-node-id", // Example self-node ID
	}

	// Provide a 40-character hex string (which is 20 bytes) for the KademliaID
	contact := &Contact{
		ID:      NewKademliaID("7461726765742d6e6f64652d69646d6f636b696461"), // 40-character hex string
		Address: "127.0.0.1",                                                 // Use Address instead of IP
	}

	// Call the real SendPingMessage function and check for proper behavior
	network.SendPingMessage(contact)
}
