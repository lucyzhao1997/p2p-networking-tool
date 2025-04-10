package main

import (
	"github.com/lucyzhao1997/p2p-networking-tool/config"
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
	log.Printf("Local server starts：%s\n", config.AppPort)
	http.ListenAndServe(config.AppPort, nil)
}