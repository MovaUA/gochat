package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8081", "Chat server address.")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var h = newHub()

func main() {
	flag.Parse()

	http.HandleFunc("/", handler)

	go h.run()

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func handler(w http.ResponseWriter, req *http.Request) {
	user, pwd, ok := req.BasicAuth()

	w.Header().Set("WWW-Authenticate", `Basic realm="chat"`)

	if !ok {
		http.Error(w, "Not authorized", 401)
		return
	}

	if pwd != "secret" {
		http.Error(w, "Not authorized", 401)
		log.Printf("user %s: invalid password\n", user)
		return
	}

	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(user, h, conn)

	h.register <- client

	go client.listen()
}
