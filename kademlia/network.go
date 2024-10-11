package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// Message struct representing the message structure used for communication
type Message struct {
	MessageType string
	Content     string
	Sender      Contact
}

type Network struct {
	Self         *Contact      //The node in wihch the network is based. This is made a pointer due to
	RoutingTable *RoutingTable //The Routing table for this network.
}

func (network *Network) Listen(ip string, port int) error { //We need to return a network
	fmt.Printf("Listening IP %s\n", ip)
	fmt.Printf("Listening port %d\n", port)

	address := fmt.Sprintf("%s:%d", ip, port)
	fmt.Printf("Listening addres %s\n", address)

	listener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	})

	if err != nil {
		fmt.Printf("ERRORRR %s\n", err)

		return err
	}
	defer listener.Close()

	fmt.Printf("Listening on %s\n", address)

	for {
		data := make([]byte, 1024) //We slice up in 1024 beacuse its relativly big and makes it possible to focus on *this* part of the bytestream which will be much smaller.
		len, remote, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}
		response, err := network.HandleMessage(data[:len])
		if err != nil {
			fmt.Println("Error when handling Message:", err)
			continue
		}
		listener.WriteToUDP(response, remote)
	}
}

<<<<<<< Updated upstream
=======
// SendPingMessage sends a PING message to a target contact to check if it's alive
func (network *Network) SendPingMessage(contact *Contact) {

	reciverID := contact.ID
	// Create the PING message
	ping := Message{
		MessageType: "PING",
		Content:     network.RoutingTable.me.Address, // The current node's address
		Sender:      network.RoutingTable.me,
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

	if msg.MessageType != "PING_ACK" {
		fmt.Println("PING_ACK not recived, we recived: ", msg.MessageType)
		return
	}
}

>>>>>>> Stashed changes
// SendMessage simulates sending a message to the given contact
func (network *Network) SendMessage(msg Message, contact *Contact) ([]byte, error) {
	connection, err := net.Dial("udp", contact.Address) //setup connection with spicific adress

	if err != nil { //Any errors during connection phase?
		return nil, err
	}

	//If no errors we move on to saving the data
	byteStream, err := json.Marshal(msg)  //Pointer due to us not wanting to send the object directly
	_, err = connection.Write(byteStream) //This row is the row that does the magic, it sends the byteStream over the network (write -> udp)

	if err != nil {
		fmt.Println("Error sending data (UDP)")
		return nil, err
	}

	// Set a deadline for the connection. If the following read operation does not complete in time, it will fail.
	deadline := time.Now().Add(15 * time.Second)
	connection.SetDeadline(deadline)

	response := make([]byte, 1024) //1024 is a balanced aproached for the amount of bytes in one slice.
	//slicing the stream up in parts is good to narrow down relevant information.
	bytesRead, err := connection.Read(response)
	if err != nil {
		fmt.Println("No response from connected node...")
		return nil, err
	}
	return response[:bytesRead], nil
}

// SendPingMessage sends a PING message to a target contact to check if it's alive
func (network *Network) SendPingMessage(contact *Contact) {
	
	reciverID := contact.ID
	ping := Message{
		MessageType: "PING",
		Content:     network.RoutingTable.me.Address, // The current node's address
		Sender:      network.RoutingTable.me,
	}

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

	if msg.MessageType != "PING_ACK" {
		fmt.Println("PING_ACK not recived, we recived: ", msg.MessageType)
		return
	}
}

// What we require Three cases:
// 1. Joining works.
// 2. Joining does not work
// 3. The contact is not reachable/does not respond.

func (network *Network) SendJoinMessage(contact *Contact) {
	reciverID := contact.ID.String()
	join := Message{
		MessageType: "JOIN",
		Content:     network.RoutingTable.me.Address,
		Sender:      network.RoutingTable.me,
	}

	response, err := network.SendMessage(join, contact)
	if err != nil {
		fmt.Printf("Failed to send JOIN message to node %s: %v", reciverID, err) 
	}

	var msg Message
	err = json.Unmarshal(response, &msg)
	if err != nil {
		fmt.Printf("No response from node %s: %v\n", reciverID, err)
	}

	if msg.MessageType != "JOIN_ACK" {
		fmt.Printf("Node %s failed to join the network %s", reciverID, join.Content)
	}
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte, contact *Contact) {
	reciverID := contact.ID.String()
	store := Message{
		MessageType: "STORE",
		Content:     network.RoutingTable.me.Address,
		Sender:      network.RoutingTable.me,
	}

	response, err := network.SendMessage(store, contact)
	if err != nil {
		fmt.Printf("Failed to send STORE message to node %s: %v", reciverID, err) 
	}

	var msg Message
	err = json.Unmarshal(response, &msg)
	if err != nil {
		fmt.Printf("No response from node %s: %v\n", reciverID, err)
	}

	if msg.MessageType != "STORE_ACK" {
		fmt.Printf("Node %s failed to store the node %s", store.Sender, store.Content)
	}

}
