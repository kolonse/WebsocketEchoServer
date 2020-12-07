// WebsocketServer project main.go
package main

import (
	"flag"

	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	PCM = 0
	LOG = 1
)

var addr = flag.String("addr", "0.0.0.0:9553", "http service address")
var useHttps = flag.Bool("https", false, "use https")
var crt = flag.String("crt", "server.crt", "https crt")
var key = flag.String("key", "server.key", "https key")

var ConnManager map[string]*websocket.Conn

//var upgrader = websocket.Upgrader{}
func serveHome(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 0, 0)
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		defer conn.Close()

		for {
			t, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("conn error: ", err)
				break
			}
			conn.WriteMessage(t, message)
		}
	}()
}

func main() {
	flag.Parse()
	http.HandleFunc("/", serveHome)
	log.Println("start listen:", *addr)

	if *useHttps {
		err := http.ListenAndServeTLS(*addr, *crt, *key, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	} else {
		err := http.ListenAndServe(*addr, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}
