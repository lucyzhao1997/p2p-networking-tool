// Handles NAT traversal logic

package nat

import (
    "fmt"
    "github.com/pion/stun"
)

//string - public ip addr
//int - public port num
//error - error message
func GetPublicIP() (string, int, error) {
	//init new client 
    c, err := stun.NewClient()
	//if error - return empty ip addr and 0 port num
    if err != nil {
        return "", 0, err
    }
    defer c.Close()
	//connects to google stun server
    if err := c.Dial("stun.l.google.com:19302"); err != nil {
        return "", 0, err
    }

	//ip and port declaration
    var publicIP string
    var publicPort int

	//Creates a new STUN message with a unique transaction ID and the "BindingRequest" type. 
	//This is the request used to get the public IP and port from the STUN server.
    message := stun.MustBuild(stun.TransactionID, stun.BindingRequest)
    
	//sends the STUN message and provides a callack function with server's response
	if err := c.Do(message, func(res stun.Event) {
        if res.Error != nil {
            panic(res.Error)
        }

		//xoraddr is used to retrieve ip and port
        var xorAddr stun.XORMappedAddress
		//extracts from res.message 
        if err := xorAddr.GetFrom(res.Message); err != nil {
            panic(err)
        }
		//reassign
        publicIP = xorAddr.IP.String()
        publicPort = xorAddr.Port
    }); err != nil {
        return "", 0, err
    }
	//nil here indicates no error

    return publicIP, publicPort, nil
}