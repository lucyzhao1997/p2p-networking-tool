package config

const (

	// The port that public server can access
	AppPort = ":4000"

	// ServerPort 
	ServerPort = ":8081"
	// TunnelPort
	TunnelPort = ":8082"
	// AppTargetPort 
	AppTargetPort = ":8083"

	// ServerIP 
	ServerIP = ""
	// ServerAddr 
	ServerAddr = ServerIP + ServerPort
	// TunnelAddr 
	TunnelAddr = ServerIP + TunnelPort

	// BufSize 
	BufSize = 1024
)