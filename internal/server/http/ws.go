package internalhttp

import (
	"github.com/filatkinen/socialnet/internal/rabbit"
	"github.com/filatkinen/socialnet/internal/rabbit/consumer"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type Ws struct {
	clientsMutex sync.Mutex
	clients      map[string]*websocket.Conn
	rabbitconn   map[string]*consumer.Comsumer
	log          *log.Logger
	upgrader     websocket.Upgrader
	exitChan     chan struct{}
	rabbitConf   rabbit.Config
}

func newWS(log *log.Logger, rabbitConf rabbit.Config) (*Ws, error) {

	return &Ws{
		clientsMutex: sync.Mutex{},
		clients:      make(map[string]*websocket.Conn),
		rabbitconn:   make(map[string]*consumer.Comsumer),
		upgrader:     websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024},
		log:          log,
		exitChan:     make(chan struct{}),
		rabbitConf:   rabbitConf,
	}, nil
}

func (ws *Ws) NewConnection(w http.ResponseWriter, r *http.Request, connID string) error {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	consum, err := consumer.NewConsumer(ws.rabbitConf, ws.log, connID)
	if err != nil {
		return err
	}

	ws.clientsMutex.Lock()
	ws.clients[connID] = conn
	ws.rabbitconn[connID] = consum
	ws.clientsMutex.Unlock()

	conn.SetCloseHandler(func(int, string) error {
		ws.log.Printf("remove connection, %s\n", connID)
		ws.rabbitconn[connID].Stop()
		ws.rabbitconn[connID].Close()
		ws.clientsMutex.Lock()
		delete(ws.clients, connID)
		delete(ws.rabbitconn, connID)
		ws.clientsMutex.Unlock()
		return nil
	})

	ws.log.Printf("New WS connection from %s", connID)

	go ws.Handle(connID)

	return nil
}

func (ws *Ws) Close() error {
	close(ws.exitChan)
	ticker := time.NewTicker(time.Millisecond * 200)
	defer ticker.Stop()

	count := 0
	select {
	case <-ticker.C:
		ws.log.Println("Waiting for close active connections...")
	default:
		if len(ws.clients) == 0 || count > 5 {
			break
		}
		count++
	}
	for k := range ws.clients {
		ws.clients[k].Close()
	}
	for k := range ws.rabbitconn {
		ws.rabbitconn[k].Close()
		ws.rabbitconn[k].Stop()
	}
	return nil
}

func (ws *Ws) Handle(connID string) {
	ws.rabbitconn[connID].Start(func(bytes []byte) {
		ws.clients[connID].WriteMessage(websocket.TextMessage, bytes)
	})
}
