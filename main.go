package main

import (
	"d7024e/cli"
	"d7024e/kademlia"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	fmt.Println("Starting the kademlia app...")
	ipf, err := GetOutboundIP()
	if err != nil {
		fmt.Println("Error getting IP: ", err)
		return
	}
	ip := ipf.String()
	if ip == "172.20.0.6" {
		StartBootstrapNode(ip)
	} else {
		StartNode(ip)
	}
	select {}
}

func StartBootstrapNode(ip string) {
	k, err := JoinNetworkBootstrap(ip, "8000")
	if err != nil {
		fmt.Println("Error joining network: ", err)
		return
	}
	go k.ListenActionChannel()
	time.Sleep(1 * time.Second)
	go k.Network.Listen(k)
	c := cli.NewCLI(k)
	go c.CliHandler()
}

func StartNode(ip string) {

	k, err := JoinNetwork(ip, "8000")
	if err != nil {
		fmt.Println("Error joining network: ", err)
		return
	}
	go k.ListenActionChannel()
	go k.Network.Listen(k)
	time.Sleep(1 * time.Second)
	DoLookUpOnSelf(k)
	c := cli.NewCLI(k)
	if c.CliHandler() {
		os.Exit(0)
	}
}

func JoinNetwork(ip string, port string) (*kademlia.Kademlia, error) {
	id := kademlia.NewRandomKademliaID()
	contact := kademlia.NewContact(id, ip+":"+port)
	contact.CalcDistance(id)
	routingTable := kademlia.NewRoutingTable(contact)
	bootStrapContact := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFFF0000000000000000000000000000000)"), "172.20.0.6:8000")
	routingTable.AddContact(bootStrapContact)

	conn, err := net.ListenPacket("udp", ":"+port)
	if err != nil {
		return nil, err
	}

	return kademlia.NewKademlia(routingTable, conn), nil
}

func JoinNetworkBootstrap(ip string, port string) (*kademlia.Kademlia, error) {
	bootStrapContact := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFFF0000000000000000000000000000000)"), ip+":"+port)
	bootStrapContact.CalcDistance(bootStrapContact.ID)
	routingTable := kademlia.NewRoutingTable(bootStrapContact)
	conn, err := net.ListenPacket("udp", ":"+port)
	if err != nil {
		return nil, err
	}

	return kademlia.NewKademlia(routingTable, conn), nil
}

func DoLookUpOnSelf(k *kademlia.Kademlia) {
	fmt.Println("Doing lookup on self")
	if k.RoutingTable == nil {
		fmt.Println("RoutingTable is nil, aborting lookup")
		return
	}
	_, _, _ = k.NodeLookup(&k.RoutingTable.Me, "")
}

func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

