package main

import (
	"encoding/json"
	"net"
	"strconv"
)

func PacketToByte(pkt Packet) []byte {
	message, _ := json.Marshal(pkt)
	return message
}
func ByteToPacket(message []byte) Packet {
	pkt := Packet{}
	json.Unmarshal(message, &pkt)
	return pkt
}

func GetUDPAddress(contact *Contact) net.UDPAddr {
	addr, port, _ := net.SplitHostPort(contact.Address) // "_" is a possible error message
	IPAddr := net.ParseIP(addr)
	intPort, _ := strconv.Atoi(port) //https://www.golangprograms.com/how-to-convert-string-to-integer-type-in-go.html

	UDPaddress := net.UDPAddr{
		IP:   IPAddr,
		Port: intPort,
		//Zone string // IPv6 scoped addressing zone
	}
	return UDPaddress
}
