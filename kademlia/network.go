package kademlia

import (
	"fmt"
	"encoding/json"
	"net"
	"time"
)

// Message struct representing the message structure used for communication
type Message struct {
	MessageType string
	Content string
	Sender Contact
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
	reciverID := contact.ID
	// Create the PING message
	ping := Message{
		MessageType:     "PING",
		Content:  network.RoutingTable.me.Address,      // The current node's address
		Sender: network.RoutingTable.me,              
	}

	// Send the PING message to the contact
	response, err := network.SendMessage(ping, contact)
	if err != nil {
		fmt.Printf("Failed to send PING message to node %s: %v\n", reciverID, err)
		return
	}

	// Wait for a response
	var msg Message //Create a empty object msg beacuse Unmarshal need to parse its result into something. And we expect this to be of struct Message.
	err = json.Unmarshal(response, &msg)
	if err != nil {
		fmt.Printf("No response from node %s: %v\n", reciverID, err)
		return
	}

	if msg.MessageType != "PING_ACK"{
		fmt.Println("PING_ACK not recived, we recived: ", msg.MessageType)
		return
	}	
}

// SendMessage simulates sending a message to the given contact
func (network *Network) SendMessage(msg Message, contact *Contact) ([]byte, error) {
	connection, err := net.Dial("udp", contact.Address) //setup connection with spicific adress
	
	if err != nil { //Any errors during connection phase?
		return nil, err
	}
	
	//If no errors we move on to saving the data 
	byteStream, err := json.Marshal(&msg) //Pointer due to us not wanting to send the object directly
	_,  err = connection.Write(byteStream) //This row is the row that does the magic, it sends the byteStream over the network (write -> udp)
	
	if err != nil {
		fmt.Println("Error sending data (UDP)")
		return nil, err
	}

	// Set a deadline for the connection. If the following read operation does not complete in time, it will fail.
	deadline := time.Now().Add(15*time.Second)
	connection.SetDeadline(deadline)

	response := make([]byte, 1024) //1024 is a balanced aproached for the amount of bytes in one slice. 
								   //slicing the stream up in parts is good to narrow down relevant information. 
	bytesRead, err := connection.Read(response)
	if err != nil{
		fmt.Println("No response from connected node...")
		return nil, err
	}
	return response[:bytesRead], nil
}

// What we require Three cases:
// 1. Joining works.
// 2. Joining does not work 
// 3. The contact is not reachable/does not respond.

func (network *Network) JoinNetwork(contact *Contact) error {
	// Create a JOIN message
	reciverID := contact.ID.String()

	join := Message{
		MessageType: "JOIN",
		Content:     network.RoutingTable.me.Address,
		Sender:    network.RoutingTable.me,
	}

	// Send the JOIN message to the contact
	response, err := network.SendMessage(join, contact)
	if err != nil {
		return fmt.Errorf("Failed to send JOIN message to node %s: %w", reciverID, err) //writing join acts as .self it seem like 
	}

	var msg Message
	err = json.Unmarshal(response, &msg)
	if err != nil {
		return fmt.Errorf("No response from node %s: %v\n", reciverID, err)
		
	}

	if msg.MessageType != "JOIN_ACK" {
		return fmt.Errorf("Node %s failed to join the network %s", reciverID, join.Content)
	}
	return fmt.Errorf("Something went wrong joining the network...")
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
