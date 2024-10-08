package cli

import (
	"d7024e/kademlia" //in go.mod we have a module this is what encapsulates our project and this is what is to be used for paths somehow.
	"fmt"
	"os"
)

func handlePing(arg string, kademliaInstance *kademlia.Kademlia, address string){
	if len(arg) == 0 {
		fmt.Println("Usage: PING <NodeID> 20+ chars")
		return
	}	
	//Kademlia make it possible to gain access to functions within this package.
	contact := kademlia.NewContact(kademlia.NewKademliaID(arg), address)
	kademliaInstance.Network.SendPingMessage(&contact) // Pass the NodeID to the PingCommand function
}

func handleJoin(arg string){
	if len(arg) == 0 {
		fmt.Println("Usage: JOIN <NodeID> 20+ chars")
		return
	}
	//kademliaInstance.Network.JoinNetwork(arg) // Pass the NodeID to the PingCommand function
}
func handleExit(arg string){
	if len(arg) == 0 {
		fmt.Println("Exiting node...")
		os.Exit(0)
	}
	fmt.Println("Usage: EXIT")

}