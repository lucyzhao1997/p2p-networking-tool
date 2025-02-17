// Helps devices discover each other.

// Facilitates the exchange of connection information (e.g., public IPs, ports, and keys).

package main

//fmt is the pacakge to use foreg. printf
import (
    "fmt"
    "net/http"
    "github.com/gorilla/websocket"
)

//websocket.upgrader upgrades http connection to websocket connection
//CheckOrigin sets to true means to permit all connection
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

//1. use upgrader to upgrade w and r 
//2. make sure websocker connection is closed when exiting the function
func handleConnection(w http.ResponseWriter, r *http.Request) {
    conn, _ := upgrader.Upgrade(w, r, nil)
    defer conn.Close()
	//assign msg with websocket's connection reading the message
	//_ will handle the errors
	//WriteMessage sends the response back to client
    for {
        _, msg, _ := conn.ReadMessage()
        fmt.Printf("Received: %s\n", msg)
        conn.WriteMessage(websocket.TextMessage, []byte("This is response from server!"))
    }
}

func main() {
	//function registration for hndleConnection
    http.HandleFunc("/ws", handleConnection)
	//hosts the server on port 8080 and listens for incoming connections
    http.ListenAndServe(":8080", nil)
}