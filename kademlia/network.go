package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
)

type Message struct {
	Type     string      
	SenderID *KademliaID 
	SenderIP string      
	TargetID string      
	TargetIP string      
	DataID   *KademliaID 
	Data     []byte
}

type Network struct {
	responseChan chan Response
	connection         net.PacketConn
}

type Response struct {
	Data            []byte    `json:"data"`
	ClosestContacts []Contact `json:"closest_contacts"`
	Target          *Contact  `json:"target"`
}

func NewNetwork(connection net.PacketConn) *Network {
	return &Network{make(chan Response), connection}
}

func (network *Network) Listen(kademliaInstance *Kademlia) {
	fmt.Println("Listening on port 8000")
	defer network.connection.Close()

	for {
		var buffer [8192]byte
		byteAmount, addr, err := network.connection.ReadFrom(buffer[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		var msg Message
		err = json.Unmarshal(buffer[:byteAmount], &msg)
		if err != nil {
			fmt.Println("Unmarshalling error, message:", err)
			continue
		}
		network.handleMessage(kademliaInstance, msg, addr)
	}
}

func (network *Network) handleMessage(kademliaInstance *Kademlia, msg Message, addr net.Addr) {
	switch msg.Type {
	case "PING":
		network.handlePing(kademliaInstance, msg, addr)

	case "STORE":
		network.handleStore(kademliaInstance, msg, addr)

	case "FIND_NODE":
		network.handleFindNode(kademliaInstance, msg, addr)

	case "FIND_DATA":
		network.handleFindData(kademliaInstance, msg, addr)
	}
}

func (network *Network) handlePing(kademliaInstance *Kademlia, msg Message, addr net.Addr) {
	PONG := Message{
		Type:     "PONG",
		SenderID: kademliaInstance.RoutingTable.Me.ID,
		SenderIP: kademliaInstance.RoutingTable.Me.Address,
	}
	data, _ := json.Marshal(PONG)
	_, err := network.connection.WriteTo(data, addr)
	if err != nil {
		fmt.Println("Error sending PONG:", err)
	} else {
		fmt.Println("Received PING. Added contact with ID:", msg.SenderID.String(), "and IP:", msg.SenderIP)
		action := Action{
			Action:   "UpdateRT",
			SenderId: msg.SenderID,
			SenderIp: msg.SenderIP,
		}
		kademliaInstance.ActionChannel <- action
	}
}

func (network *Network) SendPingMessage(sender *Contact, recipient *Contact) bool {
	PING := Message{
		Type:     "PING",
		SenderID: sender.ID,
		SenderIP: sender.Address,
	}

	response, err := network.SendMessage(sender, recipient, PING)
	if err != nil {
		fmt.Println("Error sending PING message:", err)
		return false
	}

	var msg Message
	err = json.Unmarshal(response, &msg)
	if err != nil {
		fmt.Println("Unmarshalling error, message:", err)
		return false
	}

	if msg.Type == "PONG" {
		fmt.Println("PONG from", recipient.Address)
		return true
	} else {
		fmt.Println("Unexpected message:", msg)
		return false
	}
}

func (network *Network) handleStore(kademliaInstance *Kademlia, msg Message, addr net.Addr) {
	STORE_ACK := Message{
		Type:     "STORE_ACK",
		SenderID: kademliaInstance.RoutingTable.Me.ID,
		SenderIP: kademliaInstance.RoutingTable.Me.Address,
	}
	data, _ := json.Marshal(STORE_ACK)
	_, err := network.connection.WriteTo(data, addr)
	if err != nil {
		fmt.Println("Error sending STORE_ACK:", err)
	} else {
		fmt.Println("Received STORE. Added to routing table ID:", msg.SenderID.String(), "with IP:", msg.SenderIP)
		action := Action{
			Action:   "Store",
			Hash:     msg.DataID.String(),
			Data:     msg.Data,
			SenderId: msg.SenderID,
			SenderIp: msg.SenderIP,
		}
		kademliaInstance.ActionChannel <- action
	}
}

func (network *Network) SendStoreMessage(sender *Contact, receiver *Contact, dataID *KademliaID, data []byte) bool {
	STORE := Message{
		Type:     "STORE",
		SenderID: sender.ID,
		SenderIP: sender.Address,
		DataID:   dataID,
		Data:     data,
	}

	response, err := network.SendMessage(sender, receiver, STORE)
	if err != nil {
		fmt.Println("failed to send STORE message:", err)
		return false
	}

	var STORE_ACK Message
	err = json.Unmarshal(response, &STORE_ACK)
	if err != nil {
		fmt.Println("Unmarshalling error, message:", err)
		return false
	}
	fmt.Println("Response message:", STORE_ACK.Type)
	if STORE_ACK.Type == "STORE_ACK" {
		fmt.Println("STORE_ACK from", receiver.Address)
		return true
	} else {
		fmt.Println("Unexpected message:", STORE_ACK)
		return false
	}
}

func (network *Network) handleFindData(kademliaInstance *Kademlia, msg Message, addr net.Addr) {
	if network.SendPingMessage(&kademliaInstance.RoutingTable.Me, &Contact{ID: msg.SenderID, Address: msg.SenderIP}) {
		action := Action{
			Action:   "UpdateRT",
			SenderId: msg.SenderID,
			SenderIp: msg.SenderIP,
		}
		kademliaInstance.ActionChannel <- action
	}
	action := Action{
		Action:   "LookupData",
		SenderId: msg.SenderID,
		SenderIp: msg.SenderIP,
		Hash:     msg.TargetID,
	}
	kademliaInstance.ActionChannel <- action
	responseChannel := <-network.responseChan

	response := Response{
		Data:            responseChannel.Data,
		ClosestContacts: responseChannel.ClosestContacts,
	}
	responseChannel.Data, _ = json.Marshal(response)
	_, err := network.connection.WriteTo(responseChannel.Data, addr)
	if err != nil {
		fmt.Println("Error handle closest contacts:", err)
	}
}

func (network *Network) SendFindDataMessage(sender *Contact, receiver *Contact, hash string) ([]Contact, []byte, error) {
	FINDDATA := Message{
		Type:     "FIND_DATA",
		SenderID: sender.ID,
		SenderIP: sender.Address,
		TargetID: hash,
	}

	response, err := network.SendMessage(sender, receiver, FINDDATA)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send FIND_DATA message: %v", err)
	}
	type Response struct {
		Data            []byte    `json:"data"`
		ClosestContacts []Contact `json:"closest_contacts"`
	}

	var result Response
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, nil, fmt.Errorf("Unmarshalling error, message: %v", err)
	}
	data := result.Data
	closestContacts := result.ClosestContacts

	return closestContacts, data, nil
}


func (network *Network) handleFindNode(kademliaInstance *Kademlia, msg Message, addr net.Addr) {
	fmt.Println("Received FIND_NODE")
	if network.SendPingMessage(&kademliaInstance.RoutingTable.Me, &Contact{ID: msg.SenderID, Address: msg.SenderIP}) {
		action := Action{
			Action:   "UpdateRT",
			SenderId: msg.SenderID,
			SenderIp: msg.SenderIP,
		}
		kademliaInstance.ActionChannel <- action
	} else {
		fmt.Println("Error receiving PONG")
	}
	contact := Contact{ID: NewKademliaID(msg.TargetID), Address: msg.SenderIP}
	action := Action{
		Action:   "LookupContact",
		SenderId: NewKademliaID(msg.SenderID.String()),
		SenderIp: msg.SenderIP,
		Target:   &contact,
	}
	kademliaInstance.ActionChannel <- action
	responseChannel := <-network.responseChan
	response := Response{
		Data:            responseChannel.Data,
		ClosestContacts: responseChannel.ClosestContacts,
	}
	responseChannel.Data, _ = json.Marshal(response)
	_, err := network.connection.WriteTo(responseChannel.Data, addr)
	if err != nil {
		fmt.Println("Error handle closest contacts:", err)
	}
}

func (network *Network) SendFindContactMessage(sender *Contact, receiver *Contact, target *Contact) ([]Contact, error) {
	FINDMESSAGE := Message{
		Type:     "FIND_NODE",
		SenderID: sender.ID,
		SenderIP: sender.Address,
		TargetID: target.ID.String(),
		TargetIP: target.Address,
	}

	response, err := network.SendMessage(sender, receiver, FINDMESSAGE)
	if err != nil {
		return nil, fmt.Errorf("failed to send FIND_NODE message: %v", err)
	}

	var result Response
	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, fmt.Errorf("Unmarshalling error, contacts: %v", err)
	}
	closestContacts := result.ClosestContacts
	fmt.Println("Found", len(closestContacts), "closest contacts.")
	return closestContacts, nil
}

func (network *Network) SendMessage(sender *Contact, receiver *Contact, msg interface{}) ([]byte, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", receiver.Address)
	if err != nil {
		return nil, fmt.Errorf("UDP address error: %v", err)
	}

	connection, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, fmt.Errorf("UDP error: %v", err)
	}
	defer connection.Close()

	data, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("error serializing message: %v", err)
	}

	_, err = connection.Write(data)
	if err != nil {
		return nil, fmt.Errorf("send message error: %v", err)
	}

	var buffer [8192]byte
	byteAmount, _, err := connection.ReadFromUDP(buffer[0:])
	if err != nil {
		return nil, fmt.Errorf("receiving response error: %v", err)
	}

	return buffer[:byteAmount], nil
}