package kademlia

import "fmt"

// Message struct representing the message structure used for communication
type Message struct {
	Type     string
	SenderID string
	TargetID string
}

type Network struct {
	SelfID string // Unique ID of the current node
}

func Listen(ip string, port int) {
	// TODO
}

// SendPingMessage sends a PING message to a target contact to check if it's alive
func (network *Network) SendPingMessage(contact *Contact) {
	// Create the PING message
	pingMessage := Message{
		Type:     "PING",
		SenderID: network.SelfID,      // The current node's ID
		TargetID: contact.ID.String(), // Use contact.ID.String() to get the node ID
	}

	// Send the PING message to the contact
	err := network.SendMessage(contact, pingMessage)
	if err != nil {
		fmt.Printf("Failed to send PING message to node %s: %v\n", contact.ID.String(), err)
		return
	}

	// Wait for a response (e.g., a PONG message)
	response, err := network.ReceiveResponse(contact, "PONG")
	if err != nil {
		fmt.Printf("No response from node %s: %v\n", contact.ID.String(), err)
		return
	}

	// Check if the response is a valid PONG message
	if response.Type == "PONG" {
		fmt.Printf("Node %s is alive and responding.\n", contact.ID.String())
	} else {
		fmt.Printf("Unexpected response from node %s: %v\n", contact.ID.String(), response)
	}
}

// SendMessage simulates sending a message to the given contact
func (network *Network) SendMessage(contact *Contact, msg Message) error {
	// Simulate sending a message (you can implement real logic later)
	fmt.Printf("Sending message of type '%s' from %s to %s at %s\n", msg.Type, network.SelfID, contact.ID.String(), contact.Address)
	return nil
}

// ReceiveResponse simulates receiving a response from a contact (expecting a specific message type)
func (network *Network) ReceiveResponse(contact *Contact, expectedType string) (Message, error) {
	// Simulate receiving a PONG message (you can implement real logic later)
	if expectedType == "PONG" {
		return Message{
			Type: "PONG",
		}, nil
	}
	return Message{}, fmt.Errorf("unexpected response type")
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
