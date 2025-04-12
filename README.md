# p2p-networking-tool
Simple NAT traversal 
Steps from the user's perspective:
1. User send a request via public IP address targeting the server's public endpoint
2. The servers pulls to listen on the public address, it forwards a "New Connection" notifiation to the client server when it receives
3. The client Continuously pulls from the server. it connects to the tunnel server (relay server) once it got the "New Connection"
4. It also connects to the internal service that the user wants to access
5. The client copies data between:

    Intranet service ↔ Tunnel server (e.g., io.Copy(tunnelConn, clientConn)).

    Tunnel server ↔ User (e.g., io.Copy(clientConn, tunnelConn)).
6. The server listens to the tunnel server and fetch the forwarded data
    
    via io.Copy(appConn, tunnelConn) and io.Copy(tunnelConn, appConn).
7. User recieves data from the server