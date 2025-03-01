package main

import (
	"fmt"
	"log"
	"net/url"
	"time"
	"github.com/gorilla/websocket"
)

// connectPeer simulates a WebSocket peer connection
func connectPeer(peerID string, messages chan string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Send peer ID to register
	err = conn.WriteMessage(websocket.TextMessage, []byte(peerID))
	if err != nil {
		return nil, err
	}

	// Listen for incoming messages
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("[%s] Connection closed: %v\n", peerID, err)
				return
			}
			messages <- fmt.Sprintf("[%s] Received: %s", peerID, string(msg))
		}
	}()

	return conn, nil
}

func main() {
	// Channels for receiving messages
	peerACh := make(chan string)
	peerBCh := make(chan string)

	// Connect Peer A and Peer B
	connA, err := connectPeer("PeerA", peerACh)
	if err != nil {
		log.Fatal("Failed to connect PeerA:", err)
	}
	defer connA.Close()

	connB, err := connectPeer("PeerB", peerBCh)
	if err != nil {
		log.Fatal("Failed to connect PeerB:", err)
	}
	defer connB.Close()

	// Wait for peers to establish connection
	time.Sleep(1 * time.Second)

	// Peer A sends a message to Peer B
	message := "PeerB:Hello, this is Peer A!"
	err = connA.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Fatal("Failed to send message:", err)
	}

	// Listen for responses
	select {
	case msg := <-peerACh:
		fmt.Println("✅", msg) // Expect: "Message delivered"
	case <-time.After(2 * time.Second):
		fmt.Println("❌ Peer A did not receive acknowledgment")
	}

	select {
	case msg := <-peerBCh:
		fmt.Println("✅", msg) // Expect: "From PeerA: Hello, this is Peer A!"
	case <-time.After(2 * time.Second):
		fmt.Println("❌ Peer B did not receive the message")
	}
}
