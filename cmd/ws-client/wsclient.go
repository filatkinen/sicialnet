package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

func main() {
	var serverURL string
	var token string
	flag.StringVar(&serverURL, "url", "ws://localhost:8800/post/feed/posted", "URL for ws")
	flag.StringVar(&token, "token", "4BACPA66EHIY5LJTCREDLZR2P4", "token for auth")
	flag.Parse()

	wsDialer := websocket.Dialer{}
	header := http.Header{}
	header.Add("Authorization", "Bearer "+token)
	conn, _, err := wsDialer.Dial(serverURL, header)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	signalExit := make(chan struct{})

	go func() {
		defer func() { signalExit <- struct{}{} }()
		for {
			_, buff, err := conn.ReadMessage()
			if err != nil {
				if conn != nil {
					conn.Close()
					return
				}
			}
			log.Println(string(buff))
		}
	}()

	select {
	case <-signalExit:
	case <-signalCh:
		if conn != nil {
			conn.Close()
		}
	}
}
