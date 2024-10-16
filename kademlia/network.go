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
	Self         *Contact      // The node in which the network is based. This is made a pointer to allow modifications.
	RoutingTable *RoutingTable // The Routing table for this network.
}

// Listen starts a UDP server to listen for incoming messages.
func (network *Network) Listen(ip string, port int) error {
	fmt.Printf("Listening IP %s\n", ip)
	fmt.Printf("Listening port %d\n", port)

	address := fmt.Sprintf("%s:%d", ip, port)
	fmt.Printf("Listening address %s\n", address)

	listener, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	})

	if err != nil {
		fmt.Printf("ERROR %s\n", err)
		return err
	}
	defer listener.Close()

	fmt.Printf("Listening on %s\n", address)

	for {
		data := make([]byte, 1024) // Buffer for incoming data.
		length, remote, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}
		response, err := network.HandleMessage(data[:length])
		if err != nil {
			fmt.Println("Error handling message:", err)
			continue
		}
		_, err = listener.WriteToUDP(response, remote)
		if err != nil {
			fmt.Println("Error sending response:", err)
		}
	}
}

// SendMessage sends a message to the given contact.
func (network *Network) SendMessage(msg Message, contact *Contact) ([]byte, error) {
	connection, err := net.Dial("udp", contact.Address) // Setup connection with specific address
	if err != nil {
		return nil, err
	}
	defer connection.Close()

	byteStream, err := json.Marshal(msg) // Serialize the message
	if err != nil {
		return nil, err
	}

	_, err = connection.Write(byteStream) // Send the byte stream over the network
	if err != nil {
		fmt.Println("Error sending data (UDP)")
		return nil, err
	}

	// Set a deadline for the connection. If the following read operation does not complete in time, it will fail.
	deadline := time.Now().Add(15 * time.Second)
	connection.SetDeadline(deadline)

	response := make([]byte, 1024)
	bytesRead, err := connection.Read(response)
	if err != nil {
		fmt.Println("No response from connected node...")
		return nil, err
	}
	return response[:bytesRead], nil
}

// SendPingMessage sends a PING message to a target contact to check if it's alive.
func (network *Network) SendPingMessage(contact *Contact) {
	fmt.Printf("Attempting to send PING to node %s\n", contact.ID)

	reciverID := contact.ID
	ping := Message{
		MessageType: "PING",
		Content:     network.RoutingTable.me.Address,
		Sender:      network.RoutingTable.me,
	}

	response, err := network.SendMessage(ping, contact)
	if err != nil {
		fmt.Printf("Failed to send PING message to node %s: %v\n", reciverID, err)
		return
	}

	fmt.Printf("PING message sent to node %s, waiting for response...\n", contact.ID)

	var msg Message
	err = json.Unmarshal(response, &msg)
	if err != nil {
		fmt.Printf("No response from node %s: %v\n", reciverID, err)
		return
	}

	if msg.MessageType != "PING_ACK" {
		fmt.Println("PING_ACK not received, we received:", msg.MessageType)
		return
	}

	fmt.Printf("PING_ACK received from node %s\n", contact.ID)
}

// SendJoinMessage sends a JOIN message to another node to join the network.
func (network *Network) SendJoinMessage(contact *Contact) {
	receiverID := contact.ID.String()
	join := Message{
		MessageType: "JOIN",
		Content:     network.RoutingTable.me.Address,
		Sender:      network.RoutingTable.me,
	}

	response, err := network.SendMessage(join, contact)
	if err != nil {
		fmt.Printf("Failed to send JOIN message to node %s: %v\n", receiverID, err)
		return
	}

	var msg Message
	err = json.Unmarshal(response, &msg)
	if err != nil {
		fmt.Printf("No response from node %s: %v\n", receiverID, err)
		return
	}

	if msg.MessageType != "JOIN_ACK" {
		fmt.Printf("Node %s failed to join the network %s\n", receiverID, join.Content)
		return
	}
	fmt.Printf("Successfully joined the network through node %s\n", receiverID)
}

// SendStoreMessage sends a STORE message to another node to store data.
func (network *Network) SendStoreMessage(data []byte, contact *Contact) {
	receiverID := contact.ID.String()
	store := Message{
		MessageType: "STORE",
		Content:     string(data),
		Sender:      network.RoutingTable.me,
	}

	response, err := network.SendMessage(store, contact)
	if err != nil {
		fmt.Printf("Failed to send STORE message to node %s: %v\n", receiverID, err)
		return
	}

	var msg Message
	err = json.Unmarshal(response, &msg)
	if err != nil {
		fmt.Printf("No response from node %s: %v\n", receiverID, err)
		return
	}

	if msg.MessageType != "STORE_ACK" {
		fmt.Printf("Node %s failed to store the data %s\n", receiverID, store.Content)
		return
	}
	fmt.Printf("Data successfully stored on node %s\n", receiverID)
}
