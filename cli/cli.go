package cli

import (
	"bufio"
	"d7024e/kademlia" //in go.mod we have a module this is what encapsulates our project and this is what is to be used for paths somehow.
	"fmt"
	"os"
	"strings"
)

func CommandLineInterface(kademliaInstance *kademlia.Kademlia, address string) {
	
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
		case "EXIT":
			if len(arg) == 0 {
				fmt.Println("Exiting node...")
				os.Exit(0)
			}
			fmt.Println("Usage: EXIT")
			
		default: 
			fmt.Print("Entered something bad...")

		}
	}
}
