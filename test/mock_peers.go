package main

import (
    "encoding/json"
    "net/http"
    "net/url"
    "sync"
    "testing"
    "time"

    "github.com/gorilla/websocket"
)

type Message struct {
    SenderID  string `json:"sender"`
    Recipient string `json:"recipient"`
    Content   string `json:"content"`
}

// PeerMock simulates a peer client.
type PeerMock struct {
    id   string
    conn *websocket.Conn
    mu   sync.Mutex
    lastMessage string
}

// ConnectPeer connects a mock peer to the WebSocket server.
func ConnectPeer(t *testing.T, peerID string) *PeerMock {
    serverURL := "ws://localhost:8080/ws"

    u, err := url.Parse(serverURL)
    if err != nil {
        t.Fatalf("Failed to parse server URL: %v", err)
    }

    conn, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{})
    if err != nil {
        t.Fatalf("Failed to connect peer %s: %v", peerID, err)
    }

    peer := &PeerMock{id: peerID, conn: conn}

    // Listen for incoming messages
    go func() {
        for {
            _, msg, err := conn.ReadMessage()
            if err != nil {
                break // Stop reading on error
            }
            peer.mu.Lock()
            peer.lastMessage = string(msg)
            peer.mu.Unlock()
        }
    }()

    return peer
}

// SendMessage sends a structured JSON message from one peer to another.
func (p *PeerMock) SendMessage(t *testing.T, recipient, content string) {
    msg := Message{
        SenderID:  p.id,
        Recipient: recipient,
        Content:   content,
    }
    jsonMsg, err := json.Marshal(msg)
    if err != nil {
        t.Fatalf("Failed to marshal JSON message: %v", err)
    }

    p.mu.Lock()
    defer p.mu.Unlock()
    if err := p.conn.WriteMessage(websocket.TextMessage, jsonMsg); err != nil {
        t.Fatalf("Failed to send message: %v", err)
    }
}

// WaitForMessage waits for a message to arrive within a timeout.
func (p *PeerMock) WaitForMessage(t *testing.T, expectedContent string, timeout time.Duration) bool {
    deadline := time.Now().Add(timeout)
    for time.Now().Before(deadline) {
        p.mu.Lock()
        if p.lastMessage == expectedContent {
            p.mu.Unlock()
            return true
        }
        p.mu.Unlock()
        time.Sleep(100 * time.Millisecond) // Wait and check again
    }
    return false
}

func main(t *testing.T) {
    // Step 1: Start two mock peers
    peerA := ConnectPeer(t, "PeerA")
    defer peerA.conn.Close()

    peerB := ConnectPeer(t, "PeerB")
    defer peerB.conn.Close()

    // Step 2: Peer A sends a message to Peer B
    peerA.SendMessage(t, "PeerB", "Hello, PeerB!")

    // Step 3: Wait for Peer B to receive the message
    if !peerB.WaitForMessage(t, "Hello, PeerB!", 3*time.Second) {
        t.Errorf("❌ Peer B did not receive the expected message.")
    } else {
        t.Logf("✅ Peer B received the expected message.")
    }
}
