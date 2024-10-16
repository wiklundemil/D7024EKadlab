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
			handlePing(arg, kademliaInstance, address)
		case "JOIN":
			handleJoin(arg, kademliaInstance, address)
		case "PUT":
			handlePut(arg, kademliaInstance)
		case "EXIT":
			handleExit(arg)

		default:
			fmt.Print("Entered something bad...")

		}
	}
}
