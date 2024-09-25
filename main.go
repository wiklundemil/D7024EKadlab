package main

import (
	"d7024e/kademlia" //in go.mod we have a module this is what encapsulate our project and this is what is to be used for paths somehow.
	"fmt"
	"os"
	"strings"
	"bufio"
)

func commandLineInterface(kademliaInstance *kademlia.Kademlia){
	for {
		scanner := bufio.NewReader(os.Stdin)
		fmt.Print("[Command] [INPUT] ... [INPUT]")
		fmt.Print(">>>")

		// Single quotes ('') are for runes (a single character).
		// Double quotes ("") are for strings (a sequence of characters).
		input, err := scanner.ReadString('\n') //We need to read the input when the user press ENTER ('\n')
		if err != nil {
			fmt.Print("Error while reading input...")
		}

		//When taking input we will get unnecessary data (whitespaces more then one in length etc). strings.TrimSpace removes stuff like (\t, \n, \r etc) 
		//We want a single long string of characters ([command][input]...[input])
		input = strings.TrimSpace(input)

		//As we know that we have divider with " " a single withe space we can divide the input into slices.
		slices := strings.SplitN(input, " ", 2)
		command := slices[0]
		
		//We need a way of checking if we have multiple slices (input for the commands)
		var arg string
		if len(slices) > 1{
			arg = slices[1] //as a basecase we always set the argument to be the seccond inputed value.
		}

		fmt.Print("Commands: %s, arg: %s\n", command, arg)

	}
}

func JoinNetwork(address string) *kademlia.Kademlia {
	// Create self contact
	id := kademlia.NewRandomKademliaID()
	me := kademlia.NewContact(id, address)

	// Create routing table with self as contact
	routingTable := kademlia.NewRoutingTable(me)

	// Add bootstrap contact
	bootStrapContact := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFFF0000000000000000000000000000000"), "172.20.0.6:8000")
	routingTable.AddContact(bootStrapContact)

	// Create data storage
	data := make(map[string][]byte)
	network := &kademlia.Network{}

	// Create Kademlia instance as an object
	kademliaInstance := &kademlia.Kademlia{
		RoutingTable: routingTable,
		Network:      network,
		Data:         &data,
	}

	fmt.Printf("%+v\n", kademliaInstance)
	
	return kademliaInstance
}

func main() {
	var NETWORK_IP string = "0.0.0.0"
	var NETWORK_PORT string = "3000"

	fmt.Println("Running Main function...")
	kademliaInstance := JoinNetwork(NETWORK_IP + ":" + NETWORK_PORT)
	
	//Why kademlia.Listen -> in golang we specify which package the Listen function lie in. This is enough to find the function. 
	go kademlia.Listen(NETWORK_IP, NETWORK_PORT)	//We start a goroutine by writing go first. This let us run this on a different thread. Concurrency.
	go commandLineInterface(kademliaInstance)

	select{} //This is a block for the main goroutine, used to have main running indefinitely

}

