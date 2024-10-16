package kademlia

import (
	"encoding/json"
	"fmt"
)

// HandleMessage handles incoming messages and sends responses like JOIN_ACK
func (network *Network) HandleMessage(data []byte) ([]byte, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	switch msg.MessageType {
	case "JOIN":
		// Log when a JOIN request is received
		fmt.Printf("Received JOIN request from %s\n", msg.Sender.ID)

		// Add the sender to the routing table
		network.RoutingTable.AddContact(msg.Sender)

		// Send a JOIN_ACK response back
		ack := Message{
			MessageType: "JOIN_ACK",
			Content:     network.Self.Address, // Send the current node's address as acknowledgment
			Sender:      *network.Self,
		}
		return json.Marshal(ack) // Send the JOIN_ACK message back

	default:
		fmt.Printf("Unknown message type: %s\n", msg.MessageType)
		return nil, fmt.Errorf("unknown message type")
	}
}

func (network *Network) ManagePing() Message {
	response := Message{
		MessageType: "PING_ACK",
		Content:     network.RoutingTable.me.Address,
	}
	return response
}

func (network *Network) ManageJoin(recipient Contact) Message {
	network.RoutingTable.AddContact(recipient)
	response := Message{
		MessageType: "JOIN_ACK",
		Content:     network.RoutingTable.me.Address,
	}
	return response
}

func (network *Network) ManageStore(recipient Contact) Message {
	response := Message{
		MessageType: "STORE_ACK",
		Content:     network.RoutingTable.me.Address,
	}
	return response
}
