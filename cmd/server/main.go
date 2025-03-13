package main

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

// Peer structure
type Peer struct {
    ID   string
    Conn *websocket.Conn
}

// Message structure for JSON parsing
type Message struct {
    Sender    string `json:"sender"`
    Recipient string `json:"recipient"`
    Content   string `json:"content"`
}

// Map to store connected peers
var peers = make(map[string]*Peer)
var peersMutex = sync.Mutex{}

func handleConnection(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("❌ WebSocket upgrade failed:", err)
        return
    }
    defer conn.Close()

    // Read the first message to extract the peer ID
    _, firstMsg, err := conn.ReadMessage()
    if err != nil {
        log.Println("❌ Failed to read peer ID:", err)
        return
    }

    // Parse the JSON message
    var message Message
    err = json.Unmarshal(firstMsg, &message)
    if err != nil {
        log.Println("❌ Failed to parse JSON message:", err)
        conn.WriteMessage(websocket.TextMessage, []byte("Invalid JSON format"))
        return
    }

    peerID := message.Sender // Extract sender as the peer ID

    // Store the peer
    peersMutex.Lock()
    peer := &Peer{ID: peerID, Conn: conn}
    peers[peerID] = peer
    peersMutex.Unlock()

    log.Printf("✅ Peer connected: %s", peerID)

    defer func() {
        peersMutex.Lock()
        delete(peers, peerID)
        peersMutex.Unlock()
        log.Printf("❌ Peer %s disconnected", peerID)
    }()

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Printf("❌ Peer %s error reading message: %v", peerID, err)
            return
        }

        var incomingMsg Message
        err = json.Unmarshal(msg, &incomingMsg)
        if err != nil {
            log.Println("❌ Failed to parse incoming message:", err)
            conn.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
            continue
        }

        log.Printf("📩 Received from %s: %s", incomingMsg.Sender, incomingMsg.Content)

        // Find recipient
        peersMutex.Lock()
        recipientConn, exists := peers[incomingMsg.Recipient]
        peersMutex.Unlock()

        if exists {
            err := recipientConn.Conn.WriteMessage(websocket.TextMessage, msg)
            if err != nil {
                log.Printf("❌ Error forwarding message from %s to %s: %v", incomingMsg.Sender, incomingMsg.Recipient, err)
            } else {
                log.Printf("✅ Message delivered from %s to %s", incomingMsg.Sender, incomingMsg.Recipient)
                conn.WriteMessage(websocket.TextMessage, []byte("Message delivered"))
            }
        } else {
            log.Printf("❌ Recipient %s not online", incomingMsg.Recipient)
            conn.WriteMessage(websocket.TextMessage, []byte("Recipient not online"))
        }
    }
}

func main() {
    http.HandleFunc("/ws", handleConnection)
    log.Println("🚀 Signaling server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
