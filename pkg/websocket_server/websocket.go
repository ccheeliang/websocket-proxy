package websocketserver

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Websocket struct {
	Clients   map[*websocket.Conn]bool
	Upgrader  websocket.Upgrader
	Mutex     *sync.Mutex
	Broadcast chan []byte
}

func StartWebsocketServer() *Websocket {
	ws := &Websocket{
		Mutex:     &sync.Mutex{},
		Clients:   make(map[*websocket.Conn]bool),
		Upgrader:  websocket.Upgrader{},
		Broadcast: make(chan []byte),
	}

	go ws.handleMessages()

	return ws
}

func (ws *Websocket) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	ws.Mutex.Lock()

	// Register the new client.
	ws.Clients[conn] = true
	ws.Mutex.Unlock()

	log.Println("Total Registered Clients: ", len(ws.Clients))
	// Read messages from client in order
	// to handle close connection by client.
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			switch true {
			case websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway):
				log.Printf("Connection for client %s close", conn.RemoteAddr())
			case websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure):
				log.Printf("Connection for client %s close abnormally, Error: %s", conn.RemoteAddr(), err.Error())
			default:
				log.Printf("Unhandle error: %s", err.Error())
			}

			// Remove connection from the active viewers and break the loop
			ws.Mutex.Lock()
			delete(ws.Clients, conn)
			ws.Mutex.Unlock()
			break
		}
		ws.Broadcast <- p
	}
}

func (ws *Websocket) handleMessages() {
	for {
		select {
		case msg := <-ws.Broadcast:
			log.Println("here")
			for client := range ws.Clients {
				log.Println("broadcasting")
				err := client.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Println("Error writing message:", err)
					client.Close()
					ws.Mutex.Lock()
					delete(ws.Clients, client)
					ws.Mutex.Unlock()
				}
			}
		}
	}
}

// func (ws *Websocket) handleServerInstanceMessages(msg []byte) {
// 	// Process messages received from server instances
// 	// This could involve updating the state, broadcasting to clients, etc.
// 	// Here, we broadcast the message to connected clients
// 	ws.Broadcast <- msg
// }

func ListenAndRunWebsocket(port string) {
	log.Printf("Server list and serve at %s\n", port)

	ws := StartWebsocketServer()
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.HandleWebSocket(w, r)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil))
}
