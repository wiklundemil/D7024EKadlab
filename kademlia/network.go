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
	//Setting up a UDP to receive messagese
	//IP & Port specify where node can be reached
	address := fmt.Sprintf("%s:%d", ip, port)
	listener, err := net.ListenUDP("udp", &net.UDPAddr{ //net.ListenUDP opens UDP for incoming messages on the specified ip:port
		IP:   net.ParseIP(ip),
		Port: port,
	})
	if err != nil { //returns error if address is unavailable
		return err
	}
	defer listener.Close() //after listener is set up: 'defer listener.Close()' used to ensure that listener gets closed when done

	fmt.Printf("Listening on %s\n", address) //Now that node has a UDP, it prints out "Listening on [ip]' to see where it is reachable

	for { //runs forever to wait for messages to arrive
		data := make([]byte, 1024) //buffer created to hold incoming data(message content), buffer is sliced into 1024 bytes to be a managable chunk of data
		//node waits to receive message with 'listener.ReadFromUDP(data)', this reads incoming message into the 'data' buffer and stores info about the sender in 'remote'
		len, remote, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Println("Error reading from UDP:", err) //Error while reading, makes node print the error and continues to wait for more messages
			continue
		}
		response, err := network.HandleMessage(data[:len]) //node processes received message with 'network.HandleMessage(data[:len])', func determines the kind of received message
		if err != nil {
			fmt.Println("Error when handling Message:", err) //Error while handling message, makes node print the error and continues to listen for more messages
			continue
		}
		listener.WriteToUDP(response, remote) //After handling message node creates response based of processed info, response is then sent over the UDP back to the sender, aka replying to the call!
	}
}

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

// SendMessage simulates sending a message to the a given contact/node and getting a response back.
func (network *Network) SendMessage(msg Message, contact *Contact) ([]byte, error) {
	connection, err := net.Dial("udp", contact.Address) //'net.Dial("udp", contact.Address)' used to start processes of setting up UDP connection to the spicific adress

	if err != nil { //Any errors during connection phase?
		return nil, err
	}

	//If no errors we move on to saving the data
	byteStream, err := json.Marshal(msg)  //'json.Marshal(msg)' converts msg obj to series of bytes (transmittable over the network). Pointer due to us not wanting to send the object directly
	_, err = connection.Write(byteStream) //'connection.Write(byteStream)' used to transmitt message over UDP. This row is the row that does the magic, it sends the byteStream over the network (write -> udp)

	if err != nil {
		fmt.Println("Error sending data (UDP)")
		return nil, err
	}

	// Set a deadline for the respons/connection. If the following read operation does not complete in time, it will fail.
	deadline := time.Now().Add(15 * time.Second) //15 seconds
	connection.SetDeadline(deadline)

	response := make([]byte, 1024) //1024 is a balanced aproached for the amount of bytes in one slice.
	//slicing the stream up in parts is good to narrow down relevant information.
	bytesRead, err := connection.Read(response)
	if err != nil {
		fmt.Println("No response from connected node...")
		return nil, err
	}
	return response[:bytesRead], nil //If reading response is successfull it returns the response as only 'response[:bytesRead]' aka the number of bytes actually read.
}

// What we require Three cases:
// 1. Joining works.
// 2. Joining does not work
// 3. The contact is not reachable/does not respond.

func (network *Network) JoinNetwork(contact *Contact) {
	// Create a JOIN message
	reciverID := contact.ID.String()

	join := Message{
		MessageType: "JOIN",
		Content:     network.RoutingTable.me.Address,
		Sender:      network.RoutingTable.me,
	}

	// Send the JOIN message to the contact
	response, err := network.SendMessage(join, contact)
	if err != nil {
		fmt.Errorf("Failed to send JOIN message to node %s: %w", reciverID, err) //writing join acts as .self it seem like
	}

	var msg Message
	err = json.Unmarshal(response, &msg)
	if err != nil {
		fmt.Errorf("No response from node %s: %v\n", reciverID, err)

	}

	if msg.MessageType != "JOIN_ACK" {
		fmt.Errorf("Node %s failed to join the network %s", reciverID, join.Content)
	}
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
