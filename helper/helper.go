package helper

import (
	"log"
	"net"
	"time"
)

//once user is connected to server, this method will be called. 
// ResolveTCPAddr is to resolve a TCP address, involves ip address and port 

func CreateListen(listenAddr string) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	return tcpListener, err
}

// user side initiates connection to server side, client may send packages to the server when this is successful

func CreateConnect(serverConnectAddr string) (*net.TCPConn, error) {
	// resolve address 
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverConnectAddr)
	if err != nil {
		return nil, err
	}
	//make connection to server 
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	return tcpConn, err
}

// trying to write data every 3 seconds to test if connection has been kept alive
func KeepAlive(conn *net.TCPConn) {
	for {
		_, err := conn.Write([]byte("KeepAlive"))
		if err != nil {
			log.Printf("[KeepAlive] Error %s", err)
			return
		}
		time.Sleep(time.Second * 3)
	}
}


func GetDataFromConnection(bufSize int, conn *net.TCPConn) ([]byte, error) {
	b := make([]byte, 0)
	for {
		// read data
		data := make([]byte, bufSize)
		n, err := conn.Read(data)
		if err != nil {
			return nil, err
		}
		b = append(b, data[:n]...)
		if n < bufSize {
			break
		}
	}
	return b, nil
}