package main

import (
	"d7024e/kademlia" //in go.mod we have a module this is what encapsulates our project and this is what is to be used for paths somehow.
	"d7024e/cli" //in go.mod we have a module this is what encapsulates our project and this is what is to be used for paths somehow.
	"fmt"
	"strconv"
	"time"
)

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
	go cli.CommandLineInterface(kademliaInstance, NETWORK_ADRESS)

	// Block the main function to keep the program running
	select {} // Keeps main running indefinitely
}