package cli

import (
	"d7024e/kademlia" //in go.mod we have a module this is what encapsulates our project and this is what is to be used for paths somehow.
	"fmt"
	"os"
	"strings"
)

func handlePing(arg string, kademliaInstance *kademlia.Kademlia, address string) {
	if len(arg) == 0 {
		fmt.Println("Usage: PING <NodeID> 20+ chars")
		return
	}
	//Kademlia makes it possible to gain access to functions within this package.
	contact := kademlia.NewContact(kademlia.NewKademliaID(arg), address)
	kademliaInstance.Network.SendPingMessage(&contact) // Pass the NodeID to the PingCommand function
}

func handleJoin(arg string, kademliaInstance *kademlia.Kademlia, address string) {
	if len(arg) == 0 {
		fmt.Println("Usage: JOIN <NodeID> <Node Address>")
		return
	}

	// Parse the NodeID and Address from the argument
	slices := strings.Split(arg, " ")
	if len(slices) != 2 {
		fmt.Println("Invalid arguments. Expected format: JOIN <NodeID> <Node Address>")
		return
	}

	nodeID := slices[0]      // Extract Node ID
	nodeAddress := slices[1] // Extract the Node Address (IP:Port)

	// Create a new contact using the provided NodeID and Address
	contact := kademlia.NewContact(kademlia.NewKademliaID(nodeID), nodeAddress)

	// Send a JOIN message to the contact
	kademliaInstance.Network.SendJoinMessage(&contact)
}

func handlePut(arg string, kademliaInstance *kademlia.Kademlia) {
	if len(arg) == 0 {
		fmt.Println("Usage: PUT <NodeID> 20+ chars")
		return
	}
	data := []byte(arg)
	kademliaInstance.Store(data) // Pass the NodeID to the Store function
}

func handleExit(arg string) {
	if len(arg) == 0 {
		fmt.Println("Exiting node...")
		os.Exit(0)
	}
	fmt.Println("Usage: EXIT")
}
