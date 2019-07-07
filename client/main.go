package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8081", "Chat server address")
var username = flag.String("u", "usr1", "User")
var password = flag.String("p", "", "Password")

func main() {
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	h := http.Header{"Authorization": []string{"Basic " + base64.StdEncoding.EncodeToString([]byte(*username+":"+*password))}}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), h)
	if err != nil {
		log.Fatalln("dial:", err)
		return
	}
	defer conn.Close()

	readCh := make(chan struct{})
	go func() {
		defer close(readCh)

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			if messageType != websocket.TextMessage {
				log.Println("recv: non-text message received.")
				continue
			}

			log.Printf("recv: %s", message)
		}
	}()

	writeCh := make(chan struct{})
	go func() {
		defer close(writeCh)

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			message := scanner.Text()
			if message == "quit" {
				return
			}

			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}()

	select {
	case <-interrupt:
		conn.Close()
		return
	case <-readCh:
		conn.Close()
		return
	case <-writeCh:
		conn.Close()
		return
	}
}
