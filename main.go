package main

import (
	"bufio"
	"d7024e/kademlia" //in go.mod we have a module this is what encapsulates our project and this is what is to be used for paths somehow.
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func commandLineInterface(kademliaInstance *kademlia.Kademlia, address string) {
	
	for {
		scanner := bufio.NewReader(os.Stdin)
		fmt.Print("\n[Command] [INPUT] ... [INPUT]\n>>> ")

		// Read input from the user
		input, err := scanner.ReadString('\n')
		if err != nil {
			fmt.Println("\nError while reading input...", err)
			continue
		}

		// Clean and split input
		input = strings.TrimSpace(input)
		slices := strings.SplitN(input, " ", 2)
		command := slices[0]

		// Extract the argument (NodeID) for the command
		var arg string
		if len(slices) > 1 {
			arg = slices[1]
		}

		// Handle the PING command
		switch command {
		case "PING":
			if len(arg) == 0 {
				fmt.Println("Usage: PING <NodeID> 20+ chars")
				continue
			}
			
			//Kademlia make it possible to gain access to functions within this package.
			contact := kademlia.NewContact(kademlia.NewKademliaID(arg), address)
			kademliaInstance.Network.SendPingMessage(&contact) // Pass the NodeID to the PingCommand function
		case "JOIN":
			if len(arg) == 0 {
				fmt.Println("Usage: JOIN <NodeID> 20+ chars")
				return
			}
			//kademliaInstance.Network.JoinNetwork(arg) // Pass the NodeID to the PingCommand function
		
		default: 
			fmt.Print("Entered something bad...")

		}
	}
}

func JoinNetwork(address string) *kademlia.Kademlia {
	// Create self contact
	id := kademlia.NewRandomKademliaID()
	me := kademlia.NewContact(id, address)

	// Create routing table with self as contact
	routingTable := kademlia.NewRoutingTable(me)

	// Add bootstrap contact (hardcoded for now)
	bootStrapContact := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFFF0000000000000000000000000000000"), address)
	routingTable.AddContact(bootStrapContact)

	// Create data storage
	data := make(map[string][]byte)
	network := &kademlia.Network{
		&me,                 //We pass the pointer to me
		routingTable,        //Not a defined as a pointer here due to NewROutingTable returning a pointer object of RoutingTable
	}

	// Create and return the Kademlia instance
	kademliaInstance := &kademlia.Kademlia{
		RoutingTable: routingTable,
		Network:      network,
		Data:         &data,
	}
	return kademliaInstance
}

func main() {
	var NETWORK_IP string = "0.0.0.0"
	var NETWORK_PORT int = 3000

	fmt.Println("\nRunning Main function...")

	NETWORK_PORT_STR := strconv.Itoa(NETWORK_PORT)
	var NETWORK_ADRESS string = NETWORK_IP + ":" + NETWORK_PORT_STR

	kademliaInstance := JoinNetwork(NETWORK_ADRESS)

	// Start listening on the network in a goroutine (concurrently)
	fmt.Print("111111")
	go kademliaInstance.Network.Listen(NETWORK_IP, NETWORK_PORT)
	fmt.Print("222222")

	// Ensure the listener is up before starting the CLI
	time.Sleep(1 * time.Second)

	// Start the command-line interface
	go commandLineInterface(kademliaInstance, NETWORK_ADRESS)

	// Block the main function to keep the program running
	select {} // Keeps main running indefinitely
}
