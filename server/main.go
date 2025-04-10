package main

import (
	"github.com/lucyzhao1997/p2p-networking-tool/config"
	"github.com/lucyzhao1997/p2p-networking-tool/helper"
	"io"
	"log"
	"net"
	"sync"
)

// serverConn 
var serverConn *net.TCPConn

// appConn
var appConn *net.TCPConn

// wg wait for all goroutines to finish
var wg sync.WaitGroup


func main() {
	go serverListen()
	//target server
	go appListen()
	go tunnelListen()

	//wait for all goroutines to finish
	wg.Add(1)
	wg.Wait()
}

func serverListen() {
	//this is for accepting client connection
	tcpListener, err := helper.CreateListen(config.ServerAddr)
	if err != nil {
		panic(err)
	}
	log.Printf("Server Address：%s\n", tcpListener.Addr().String())

	//wait to receive client connection
	for {
		serverConn, err = tcpListener.AcceptTCP()
		if err != nil {
			log.Printf("Connection failed, error log：%s\n", err.Error())
			return
		}
		go helper.KeepAlive(serverConn)
	}
}

// NAT traversal, this is for navigating between client side and server side data and connnections
func tunnelListen() {
	tcpListener, err := helper.CreateListen(config.TunnelAddr)
	if err != nil {
		panic(err)
	}
	log.Printf("tunnel server address：%s\n", tcpListener.Addr().String())
	for {
		tunnelConn, err := tcpListener.AcceptTCP()
		if err != nil {
			log.Printf("tunnel server failed, error log：%s\n", err.Error())
			return
		}
		// data transfer
		go io.Copy(appConn, tunnelConn)
		go io.Copy(tunnelConn, appConn)
	}
}

// endpoint server
func appListen() {
	//监听目的服务端
	tcpListener, err := helper.CreateListen(config.AppTargetPort)
	if err != nil {
		panic(err)
	}
	log.Printf("End point server address is: %s\n", tcpListener.Addr().String())

	for {
		appConn, err = tcpListener.AcceptTCP()
		if err != nil {
			log.Printf("end point conenction failed, error log：%s\n", err.Error())
			return
		}
		_, err := serverConn.Write([]byte("New Connection"))
		if err != nil {
			log.Printf("message send failed, error：%s\n", err.Error())
		}
	}
}