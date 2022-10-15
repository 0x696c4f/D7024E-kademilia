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
	messageFindValue := networkOther.NewPacket("find_Value")
	messageFindValue.Message = MessageBody{
		Hash: HashData("hejsan"),
	}
	b := []byte("ABC€")

	networkOther.Store(b)
	messageStore := networkOther.NewPacket("store_Value")

	var messageLocalGet = Packet{
		SendingContact: &networkMe.Node.RoutingTable.me,
		RPC:            "local_get",
		Message: MessageBody{
			TargetID: NewKademliaID(HashData("hejsan")),
		},
	}
	data := []byte("ABC€")

	var messageLocalPut = Packet{
		SendingContact: &networkMe.Node.RoutingTable.me,
		RPC:            "local_put",
		Message: MessageBody{
			Data: data,
		},
	}

	var messageRefresh = Packet{
		SendingContact: &networkMe.Node.RoutingTable.me,
		RPC:            "refresh",
		Message: MessageBody{
			TargetID: NewKademliaID(HashData("hejsan")),
		},
	}

	var messageForget = Packet{
		SendingContact: &networkMe.Node.RoutingTable.me,
		RPC:            "local_forget",
		Message: MessageBody{
			TargetID: NewKademliaID(HashData("hejsan")),
		},
	}

	messageFind.Message = MessageBody{
		TargetID: target.ID,
	}

	networkMe.MessageHandler(messageEmpty)
	networkMe.MessageHandler(messagePing)
	networkMe.MessageHandler(messageFind)
	networkMe.MessageHandler(messageFindValue)
	networkMe.MessageHandler(messageStore)

	fmt.Println("test localget")
	networkMe.MessageHandler(messageLocalGet)
	networkMe.MessageHandler(messageLocalPut)

	networkMe.MessageHandler(messageRefresh)
	networkMe.MessageHandler(messageForget)

	fmt.Println("-------------------------")
	fmt.Println("")
}
