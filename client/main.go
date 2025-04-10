// Connects to the signaling server.

// Handles NAT traversal (STUN/TURN).

// Manages WireGuard connections.

// Provides a CLI or GUI for users to interact with the tool.

package main

import (
    "github.com/lucyzhao1997/p2p-networking-tool/config"
    "github.com/lucyzhao1997/p2p-networking-tool/helper"
    "log"
    "io"
)

func main() {
    //connect to server
    conn, err := helper.CreateConnect(config.ServerAddr)
    if err != nil {
		panic(err)
	}
	log.Printf("Connected Successfully, the address is：%s\n", conn.RemoteAddr().String())
    for {
		
		data, err := helper.GetDataFromConnection(config.BufSize, conn)
		if err != nil {
			log.Printf("Failed to read data, error log：%s\n", err.Error())
			continue
		}
		log.Printf("Data recieved：%s\n", string(data))
		
		if string(data) == "New Connection" {
			//connect to tunnel server
			go messgaeForward()
		}
	}
}
func messgaeForward() {
	// connect to tunnel server
	tunnelConn, err := helper.CreateConnect(config.TunnelAddr)
	if err != nil {
		panic(err)
	}

	// connect to client side service
	clientConn, err := helper.CreateConnect(config.AppPort)
	if err != nil {
		panic(err)
	}

	go io.Copy(clientConn, tunnelConn)
	go io.Copy(tunnelConn, clientConn)
}