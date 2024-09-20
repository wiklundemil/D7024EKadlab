package kademlia

import(
	"encoding/json"
	"errors"
	"fmt"
)

func (network *Network) HandleMessage(byteStream []byte) ([]byte, error){
	var msg Message //msg of type Message
	
	err := json.Unmarshal(byteStream, &msg) //pointer beacuse we do not want to use our msg directly
	if err != nil {                         //There is something wrong with json format, something went wrong reconstructing the json
		fmt.Println("Something wrong with byteStream")
	}

	if msg.MessageType == "PING" {
		response := network.ManagePing()
		data, err := json.Marshal(response)
		return data, err
	}

	if msg.MessageType == "JOIN" {
		response := network.ManageJoin(msg.Sender)
		data, err := json.Marshal(response)
		return data, err
	}
	return nil, errors.New("Unkown command made...")
}

func (network *Network) ManagePing() Message {
	response := Message{
		MessageType: "PING_ACK",
		Content:     network.RoutingTable.me.Address,
	}
	return response
}

func (network *Network) ManageJoin(recipient Contact) Message {
	network.RoutingTable.AddContact(recipient)
	response := Message{
		MessageType: "JOIN_ACK",
		Content:     network.RoutingTable.me.Address,
	}
	return response
}