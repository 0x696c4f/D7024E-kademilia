package main

import (
	"fmt"
	"testing"
)

func TestMessageHandler(t *testing.T) {
	fmt.Println("implement messageHandler testing")
	target := NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002")

	networkMe := NewNetwork("localhost:8000")
	networkOther := NewNetwork("172.17.0.2:8000")

	messageEmpty := networkOther.NewPacket("empty")
	messagePing := networkOther.NewPacket("ping")
	messageFind := networkOther.NewPacket("find_Node")
	messageFind.Message = MessageBody{
		TargetID: target.ID,
	}

	networkMe.MessageHandler(messageEmpty)
	networkMe.MessageHandler(messagePing)
	networkMe.MessageHandler(messageFind)

	fmt.Println("-------------------------")
	fmt.Println("")
}
