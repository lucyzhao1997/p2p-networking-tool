package helper

import (
	"log"
	"net"
	"time"
)

//once user is connected to server, this method will be called. 
// ResolveTCPAddr is to resolve a TCP address, involves ip address and port 
// ListenTCP函数监听TCP地址，addr则是一个TCP地址，
// 如果addr的端口字段为0，函数将选择一个当前可用的端口，
// 返回值l是一个net.Listener接口，可以用来接收连接。
func CreateListen(listenAddr string) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	return tcpListener, err
}

// CreateConnect 连接，参数为服务端地址connectAddr，返回 TCPConn，通过 net.ResolveTCPAddr 解析地址，通过 net.DialTCP 连接服务端
// client side initiates connection to server side, client may send packages to the server when this is successful

// DialTCP函数在网络协议tcp上连接本地地址laddr和远端地址raddr，如果laddr为nil，则自动选择本地地址，如果raddr为nil，则函数在建立连接之前不会尝试解析地址，一般用于客户端。
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