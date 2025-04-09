package p2pnetworkingtool

import (
    "encoding/json"
    "log"
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
    id          string
    conn        *websocket.Conn
    mu          sync.Mutex
    lastMessage string
}

// ConnectPeer connects a mock peer to the WebSocket server.
func ConnectPeer(t *testing.T, peerID string) *PeerMock {
    serverURL := "ws://localhost:8080/ws"

    u, err := url.Parse(serverURL)
    if err != nil {
        t.Fatalf("Failed to parse server URL: %v", err)
    }

    conn, resp, err := websocket.DefaultDialer.Dial(u.String(), http.Header{})
    if err != nil {
        t.Fatalf("Failed to connect peer %s: %v, Response: %+v", peerID, err, resp)
    }

    // Send initial message with peer ID
    initMsg := Message{
        SenderID: peerID,
    }
    jsonInitMsg, err := json.Marshal(initMsg)
    if err != nil {
        t.Fatalf("Failed to marshal initial JSON message: %v", err)
    }

    if err := conn.WriteMessage(websocket.TextMessage, jsonInitMsg); err != nil {
        t.Fatalf("Failed to send initial message: %v", err)
    }

    peer := &PeerMock{id: peerID, conn: conn}

    // Listen for incoming messages
    go func() {
        for {
            _, msg, err := conn.ReadMessage()
            if err != nil {
                if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                    t.Logf("Peer %s error reading message: %v", peerID, err)
                }
                return // Stop reading on error
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
    t.Logf("Peer %s sending message to %s: %s", p.id, recipient, content)
}

// WaitForMessage waits for a message to arrive within a timeout.
func (p *PeerMock) WaitForMessage(t *testing.T, expectedContent string, timeout time.Duration) bool {
    deadline := time.Now().Add(timeout)
    for time.Now().Before(deadline) {
        p.mu.Lock()
        receivedMessage := p.lastMessage
        p.mu.Unlock()
        if receivedMessage == expectedContent {
            t.Logf("Peer %s received the expected message: %s", p.id, receivedMessage)
            return true
        }
        time.Sleep(100 * time.Millisecond) // Wait and check again
    }
    t.Logf("Peer %s did not receive the expected message. Last received: %s", p.id, p.lastMessage)
    return false
}

// TestMockPeers runs the mock peer test.
func TestMockPeers(t *testing.T) {
    // Start the server in a separate goroutine
    go func() {
        log.Println("Starting server...")
        StartServer()
    }()
    // Give the server some time to start up
    time.Sleep(2 * time.Second)

    // Step 1: Start two mock peers
    peerA := ConnectPeer(t, "PeerA")
    defer peerA.conn.Close()

    peerB := ConnectPeer(t, "PeerB")
    defer peerB.conn.Close()

    // Step 2: Peer A sends a message to Peer B
    peerA.SendMessage(t, "PeerB", "Hello, PeerB!")

    // Step 3: Wait for Peer B to receive the message
    expectedMessage := "{\"sender\":\"PeerA\",\"recipient\":\"PeerB\",\"content\":\"Hello, PeerB!\"}"
    if !peerB.WaitForMessage(t, expectedMessage, 5*time.Second) {
        t.Errorf("❌ Peer B did not receive the expected message.")
    } else {
        t.Logf("✅ Peer B received the expected message.")
    }
}