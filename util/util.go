package util

import (
	"crypto/rand"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"net"
	"os"
)

const (
	CONNECT   = iota // connect
	HEARTBEAT        // heartbeat
	S_TO_C           // The server forwards user data to the client
	C_TO_S           // The client forwards the data returned by the target service to the server
)

// LoadConfig loads the configuration from the given path
func LoadConfig[T any](path string) *T {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading config file: %v", err)
		panic(err)
	}
	var config T
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Printf("Error parsing config file: %v", err)
		panic(err)
	}

	return &config
}

func CreateUDPListen(listenAddr string, listenPort string) (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", listenAddr+":"+listenPort)
	if err != nil {
		return nil, err
	}
	udpListener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	return udpListener, err
}

func CreateUDPConnect(connectAddr string, listenPort string) (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", connectAddr+":"+listenPort)
	if err != nil {
		return nil, err
	}
	tcpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}
	return tcpConn, err
}

// GenerateUUID create SessionID. The SessionID is used to address the issue of handling multiple connections on a single port.
func GenerateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant is 10

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
