package main

import (
	"log"
	"net"
	"time"
	"util"
)

// Define configuration
type Config struct {
	Client struct {
		Host         string        `yaml:"host"`
		Port         string        `yaml:"port"`
		ClientId     string        `yaml:"clientId"`
		PackSize     int           `yaml:"packSize"`
		HoleInterval time.Duration `yaml:"holeInterval"`
	} `yaml:"client"`
}

// Configuration
var cfg *Config

var udpConn *net.UDPConn

func UDPClient() {
	var err error
	for {
		udpConn, err = util.CreateUDPConnect(cfg.Client.Host, cfg.Client.Port)
		if err != nil {
			log.Printf("UDPClient connect Error: %s\n", err.Error())
			time.Sleep(cfg.Client.HoleInterval)
			continue
		}
		log.Printf("UDPClient connect success: %s", udpConn.RemoteAddr().String())
		data := make([]byte, cfg.Client.PackSize)
		clientIdLen := len(cfg.Client.ClientId)
		data[0] = util.CONNECT
		copy(data[1:], cfg.Client.ClientId)
		_, _ = udpConn.Write(data[:clientIdLen+1])
		//go ping(udpConn) // You shouldn't need to ping, just reconnect every once in a while
		go handleConn(udpConn)
		// Every once in a while, holes are drilled into the server
		time.Sleep(cfg.Client.HoleInterval)
	}
}

func handleConn(conn *net.UDPConn) {
	buf := make([]byte, cfg.Client.PackSize)
	var n int
	var length int
	for {
		// Accept the S data first, then send it to T, accept the return of T, and then send it to S
		n, _, _ = conn.ReadFromUDP(buf)
		if n == 0 {
			continue
		}
		sessionId := string(buf[1:37])
		thostLen := buf[37]
		thost := string(buf[38 : 38+thostLen])
		tportLen := buf[38+thostLen]
		tport := string(buf[39+thostLen : 39+thostLen+tportLen])
		data := buf[39+thostLen+tportLen : n]

		tConn, _ := util.CreateUDPConnect(thost, tport)
		length, _ = tConn.Write(data)
		log.Printf("[%s] S -> T len=%d", sessionId, length)
		n, _, _ = tConn.ReadFromUDP(buf)

		cToSData := make([]byte, cfg.Client.PackSize)
		cToSData[0] = util.C_TO_S
		copy(cToSData[1:], sessionId)
		copy(cToSData[37:], buf[:n])
		length, _ = conn.Write(cToSData[:37+n])
		log.Printf("[%s] T -> S len=%d", sessionId, length-37)
	}
}

func main() {

	cfg = util.LoadConfig[Config]("./config/app.yml")

	UDPClient()

}
