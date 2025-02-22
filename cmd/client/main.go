// Connects to the signaling server.

// Handles NAT traversal (STUN/TURN).

// Manages WireGuard connections.

// Provides a CLI or GUI for users to interact with the tool.

package main

import (
    "fmt"
    "github.com/lucyzhao1997/p2p-networking-tool/internal/nat"
    "github.com/lucyzhao1997/p2p-networking-tool/internal/security"
    "github.com/lucyzhao1997/p2p-networking-tool/pkg/config"
    "log"
)

func main() {

    cfg := config.LoadConfig()

    // Generate RSA Key Pair
    privateKey, publicKey, err := security.GenerateKeyPair()
    if err != nil {
        log.Fatalf("Failed to generate key pair: %v", err)
    }
    fmt.Println("Public Key:", publicKey)

    //get public ip addr
    ip, port, err := nat.GetPublicIP()
    if err != nil {
		//handle error message
        panic(err)
    }
    fmt.Printf("Public IP: %s, Port: %d\n", ip, port)
}