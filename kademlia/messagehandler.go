package kademlia

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Takes in an incoming bytestream (raw data) from another node adn converts into a msg obj that the node can understand and respond to
func (network *Network) HandleMessage(byteStream []byte) ([]byte, error) {
	var msg Message //msg of type Message

	err := json.Unmarshal(byteStream, &msg) //Converting incoming bytestream back into a msg obj using 'json.Unmarshalpointer'. Pointer beacause we do not want to use our msg directly.
	if err != nil {                         //There is something wrong with json format, something went wrong reconstructing the json
		fmt.Println("Something wrong with byteStream")
	}

	//Checks what type of message it is (PING or JOIN)

	if msg.MessageType == "PING" {
		response := network.ManagePing() //calls ManagePing that returns a PING_ACK
		data, err := json.Marshal(response)
		return data, err
	}

	if msg.MessageType == "JOIN" {
		response := network.ManageJoin(msg.Sender) //calls ManageJoin that returns a JOIN_ACK
		data, err := json.Marshal(response)
		return data, err
	}
	return nil, errors.New("Unkown command made...") //If message type not recoginized
}

func (network *Network) ManagePing() Message {
	response := Message{ //creates a PING_ACK message in response including node’s address
		MessageType: "PING_ACK",
		Content:     network.RoutingTable.me.Address,
	}
	return response
}

func (network *Network) ManageJoin(recipient Contact) Message {
	network.RoutingTable.AddContact(recipient) //calls 'network.RoutingTable.AddContact(recipient)' to add the sender(the new node) to the Routing Table
	response := Message{                       //creates a JOIN_ACK message in response including node’s address
		MessageType: "JOIN_ACK",
		Content:     network.RoutingTable.me.Address,
	}
	return response
}
