// Helps devices discover each other.

// Facilitates the exchange of connection information (e.g., public IPs, ports, and keys).

package main

//fmt is the pacakge to use foreg. printf
import (
    "fmt"
    "net/http"
    "sync"
    "github.com/gorilla/websocket"
)

//websocket.upgrader upgrades http connection to websocket connection
//CheckOrigin sets to true means to permit all connection
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

type Peer struct {
    ID string
    Conn *websocket.Conn
}

//store active peers where key is the unique peer id
var peers = make(map[string]*Peer)

var peersMutex = sync.Mutex{}

//1. use upgrader to upgrade w and r 
//2. make sure websocker connection is closed when exiting the function
func handleConnection(w http.ResponseWriter, r *http.Request) {
    conn, _ := upgrader.Upgrade(w, r, nil)
    defer conn.Close()

    // Read the peer ID from the first message sent by the client.
    _, peerID, err := conn.ReadMessage()
    if err != nil {
        log.Println("Failed to read peer ID:", err)
        return
    }
    //find the peer to connect with
    peersMutex.Lock()
    peer := &Peer{ID: string(peerID), Conn: conn}
    peers[peer.ID] = peer

    peersMutex.Unlock()
    fmt.Println("Peer connected:", peer.ID)
	
    //when closing connection, delete the peer as well
    defer func() {
		peersMutex.Lock()
		delete(peers, peer.ID)
		peersMutex.Unlock()
		fmt.Printf("Peer %s disconnected\n", peerIDStr)
	}()
    //assign msg with websocket's connection reading the message
	//_ will handle the errors
	//WriteMessage sends the response back to client
    for {
        _, msg, _ := conn.ReadMessage()
        fmt.Printf("Received: %s\n", msg)
        // Parse the message to extract the recipient ID and the actual message.
		// Assumes the message format is "recipientID:message".
		var recipientID, message string
		n, _ := fmt.Sscanf(string(msg), "%s:%s", &recipientID, &message)
		if n != 2 {
			// If the message format is invalid, send an error message back to the sender.
			conn.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
			continue
		}

		// Check if the recipient is online (i.e., in the peers map).
		peersMutex.Lock()
		recipientConn, exists := peers[recipientID]
		peersMutex.Unlock()

		if exists {
			// If the recipient is online, forward the message to them.
			err := recipientConn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("From %s: %s", peerIDStr, message)))
			if err != nil {
				fmt.Println("Error forwarding message:", err)
			} else {
				// Send an acknowledgment back to the sender.
				conn.WriteMessage(websocket.TextMessage, []byte("Message delivered"))
			}
		} else {
			// If the recipient is not online, notify the sender.
			conn.WriteMessage(websocket.TextMessage, []byte("Recipient not online"))
		}
    }
}

func main() { 
	//function registration for handleConnection
    http.HandleFunc("/ws", handleConnection)
    // Start the HTTP server on port 8080.
	fmt.Println("Signaling server started on :8080")
	//hosts the server on port 8080 and listens for incoming connections
    http.ListenAndServe(":8080", nil)
}