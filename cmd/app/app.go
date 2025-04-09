package main

import (
	"p2p-networking-tool/cmd/config"
	"encoding/json"
	"log"
	"net/http"
)

//test NAT traversal
func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		q := request.URL.Query()
		b, err := json.Marshal(q)
		if err != nil {
			log.Println(err)
		}
		writer.Write(b)
	})
	log.Printf("Local server starts：%s\n", constant.AppPort)
	http.ListenAndServe(constant.AppPort, nil)
}