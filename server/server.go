package main

import (
	"log"
	"net"
	"sync"
	"util"
)

// Define configuration
type Config struct {
	Server struct {
		Port    string `yaml:"port"`
		Clients []struct {
			ClientID string `yaml:"clientId"`
			Port     string `yaml:"port"`
			THost    string `yaml:"tHost"`
			TPort    string `yaml:"tPort"`
		} `yaml:"clients"`
		PackSize int `yaml:"packSize"`
	} `yaml:"server"`
}

// wg Used to wait for all coroutines to finish.
var wg sync.WaitGroup

// rwLock Used to protect sessionListenerMap, sessionAddrMap, clientAddrMap.
var rwLock sync.RWMutex

// Configuration
var cfg *Config

// sessionId -> listener
var sessionListenerMap map[string]*net.UDPConn

// sessionId -> UserAddr
var sessionAddrMap map[string]*net.UDPAddr

// clientId -> ClientAddr
var clientAddrMap map[string]*net.UDPAddr

// Listen to the client's UDP packets
var clientUDPConn *net.UDPConn

// initialize
func initialize() {
	// load config
	cfg = util.LoadConfig[Config]("./config/app.yml")

	sessionListenerMap = make(map[string]*net.UDPConn)
	sessionAddrMap = make(map[string]*net.UDPAddr)
	clientAddrMap = make(map[string]*net.UDPAddr)

	var err error
	clientUDPConn, err = util.CreateUDPListen("", cfg.Server.Port)
	if err != nil {
		log.Fatalf("CreateUDPListen Error: %s\n", err.Error())
	}
	log.Printf("Server Started at %s", clientUDPConn.LocalAddr())
}

func handleConn(conn *net.UDPConn) {
	buf := make([]byte, cfg.Server.PackSize)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if n == 0 || err != nil {
			continue
		}
		if buf[0] == util.CONNECT {
			clientId := string(buf[1:n])
			for _, client := range cfg.Server.Clients {
				if client.ClientID == clientId {
					rwLock.Lock()
					clientAddrMap[client.ClientID] = addr
					log.Printf("UDPClient connect success: [%s] %s", clientId, addr.String())
					rwLock.Unlock()
				}
			}
			// 非法ClientID
		} else if buf[0] == util.C_TO_S {
			sessionId := string(buf[1:37])
			rwLock.RLock()
			length := 0
			length, _ = sessionListenerMap[sessionId].WriteToUDP(buf[37:n], sessionAddrMap[sessionId])
			log.Printf("[%s] C -> U len=%d", sessionId, length)
			delete(sessionListenerMap, sessionId)
			delete(sessionAddrMap, sessionId)
			rwLock.RUnlock()
		}
	}
}

func portListen(port string, clientID string, thost string, tport string) {
	listener, err := util.CreateUDPListen("", port)
	if err != nil {
		log.Fatalf("User listen Error: %s\n", err)
		return
	}
	log.Printf("User listen SUCCESS：%s\n", listener.LocalAddr().String())
	buf := make([]byte, cfg.Server.PackSize)
	for {
		n, userAddr, _ := listener.ReadFromUDP(buf)
		sessionID := util.GenerateUUID()

		rwLock.Lock()
		sessionListenerMap[sessionID] = listener
		sessionAddrMap[sessionID] = userAddr
		rwLock.Unlock()

		data := make([]byte, cfg.Server.PackSize)
		data[0] = util.S_TO_C
		thostLen := byte(len(thost))
		tportLen := byte(len(tport))
		copy(data[1:], sessionID)
		data[37] = thostLen
		copy(data[38:], thost)
		data[38+thostLen] = tportLen
		copy(data[39+thostLen:], tport)
		copy(data[39+thostLen+tportLen:], buf[:n])

		rwLock.RLock()
		_, _ = clientUDPConn.WriteToUDP(data[:int(39+thostLen+tportLen)+n], clientAddrMap[clientID])
		log.Printf("[%s] U -> C len=%d", sessionID, n)
		rwLock.RUnlock()
	}
}

func main() {

	initialize()

	go handleConn(clientUDPConn)

	for _, client := range cfg.Server.Clients {
		go portListen(client.Port, client.ClientID, client.THost, client.TPort)
	}

	wg.Add(1)
	wg.Wait()
}
