package kademlia

import "fmt"

// Message struct representing the message structure used for communication
type Message struct {
	MessageType string
	Content string
	TargetID string
}

type Network struct {
	Self *Contact 		       //The node in wihch the network is based. This is made a pointer due to 
	RoutingTable *RoutingTable //The Routing table for this network. 
}

func Listen(ip string, port int) {
	// TODO
}

// SendPingMessage sends a PING message to a target contact to check if it's alive
func (network *Network) SendPingMessage(contact *Contact) {
	// Create the PING message
	ping := Message{
		MessageType:     "PING",
		Content:  network.RoutingTable.me.Address,      // The current node's ID
		TargetID: contact.ID.String(), // Use contact.ID.String() to get the node ID
	}

	// Send the PING message to the contact
	err := network.SendMessage(contact, ping)
	if err != nil {
		fmt.Printf("Failed to send PING message to node %s: %v\n", ping.TargetID, err)
		return
	}

	// Wait for a response (e.g., a PONG message)
	response, err := network.ReceiveResponse(contact, ping.MessageType)
	if err != nil {
		fmt.Printf("No response from node %s: %v\n", ping.TargetID, err)
		return
	}

	// Check if the response is a valid PONG message
	if response.MessageType == ping.MessageType + "_ACK" {
		fmt.Printf("Node %s is alive and responding.\n", ping.TargetID)
	} else {
		fmt.Printf("Unexpected response from node %s: %v\n", ping.TargetID, response)
	}
}

// SendMessage simulates sending a message to the given contact
func (network *Network) SendMessage(contact *Contact, msg Message) error {
	// Simulate sending a message (you can implement real logic later)
	fmt.Printf("Sending message of type '%s' from %s to %s at %s\n", msg.MessageType, network.Self, msg.TargetID, contact.Address)
	return nil
}

// What we require Three cases:
// 1. Joining works.
// 2. Joining does not work 
// 3. The contact is not reachable/does not respond.

func (network *Network) JoinNetwork(contact *Contact) error {
	// Create a JOIN message
	join := Message{
		MessageType: "JOIN",
		Content:     network.RoutingTable.me.Address,
		TargetID:    contact.ID.String(),
	}

	// Send the JOIN message to the contact
	err := network.SendMessage(contact, join)
	if err != nil {
		return fmt.Errorf("failed to send JOIN message to node %s: %w", join.TargetID, err) //writing join acts as .self it seem like 
	}

	// Wait for a response (e.g., a JOIN_ACK message)
	response, err := network.ReceiveResponse(contact, join.MessageType)
	if err != nil {
		return fmt.Errorf("no response from node %s: %w", join.TargetID, err)
	}

	// Check if the response is a valid JOIN_ACK message
	if response.MessageType == join.MessageType + "_ACK" {
		// Add the contact to the routing table
		if network.RoutingTable == nil {
			network.RoutingTable = NewRoutingTable(*contact) // Initialize RoutingTable if not already
		}
		network.RoutingTable.AddContact(*contact)
		fmt.Printf("Node %s has successfully joined the network.\n", join.TargetID)
		return nil
	}

	return fmt.Errorf("unexpected response from node %s: %v", join.TargetID, response)
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


// ReceiveResponse simulates receiving a response from a contact with the expected message type
func (network *Network) ReceiveResponse(contact *Contact, expectedType string) (Message, error) {
	// Simulate receiving a response; in a real implementation, this would involve
	// waiting for and parsing incoming messages.
	if expectedType == "PING" {
		return Message{
			MessageType: "PING_ACK",
		}, nil
	} 
	if expectedType == "JOIN" {
		return Message{
			MessageType: "JOIN_ACK",
		}, nil
	}
	return Message{}, fmt.Errorf("unexpected response type")
}