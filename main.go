package main

import (
	"d7024e/cli"
	"d7024e/kademlia"
	"fmt"
	"os"
	"strconv"
	"time"
)

func JoinNetwork(address string) *kademlia.Kademlia {
	id := kademlia.NewRandomKademliaID()   // Generate a random node ID
	me := kademlia.NewContact(id, address) // Create a Contact for the current node

	routingTable := kademlia.NewRoutingTable(me) // Correctly initialize routing table
	network := &kademlia.Network{
		Self:         &me,          // Correctly initialize 'Self' with the current node's contact
		RoutingTable: routingTable, // Pass the initialized routing table
	}

	data := make(map[string][]byte) // Initialize the data storage map

	kademliaInstance := &kademlia.Kademlia{
		RoutingTable: routingTable,
		Network:      network,
		Data:         &data,
	}
	return kademliaInstance
}

func main() {
	var NETWORK_IP string = "0.0.0.0"
	var NETWORK_PORT int

	// Get port number from the environment variable or command-line argument
	portStr, ok := os.LookupEnv("NODE_PORT")
	if !ok {
		if len(os.Args) > 1 {
			portStr = os.Args[1]
		} else {
			fmt.Println("No valid port provided. Defaulting to port 3000.")
			portStr = "3000" // Default port if nothing is specified
		}
	}

	port, err := strconv.Atoi(portStr)
	if err != nil || port < 3000 || port > 4000 {
		fmt.Println("Invalid port number provided. Please provide a valid port number between 3000 and 4000.")
		return
	}
	NETWORK_PORT = port

	fmt.Println("\nRunning Main function...")
	fmt.Printf("Listening on %s:%d\n", NETWORK_IP, NETWORK_PORT)

	NETWORK_ADDRESS := fmt.Sprintf("%s:%d", NETWORK_IP, NETWORK_PORT)
	kademliaInstance := JoinNetwork(NETWORK_ADDRESS)

	go kademliaInstance.Network.Listen(NETWORK_IP, NETWORK_PORT)
	time.Sleep(1 * time.Second)

	go cli.CommandLineInterface(kademliaInstance, NETWORK_ADDRESS)

	// Keep the program running indefinitely
	select {}
}
