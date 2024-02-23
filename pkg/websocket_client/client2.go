package websocketclient

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var messageChannel = make(chan []byte) // Channel to send data to the WebSocket server

// wsServerURL have to be in "ws://{host}:{port}/ws" and not http scheme
func StartServer2(port string, wsServerURL string) {
	log.Printf("Server list and serve at %s\n", port)

	// Connect to the WebSocket server
	serverConn, _, err := websocket.DefaultDialer.Dial(wsServerURL, nil)
	if err != nil {
		log.Println("Error connecting to WebSocket server:", err)
		return
	}
	defer serverConn.Close()

	// Proxy messages between server instance and WebSocket server
	go broadcastToClients(serverConn)

	// r := gin.Default()
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})
	// r.GET("/broadcast", func(c *gin.Context) {
	// 	msg := c.Query("message")
	// 	messageChannel <- []byte(msg)
	// 	c.JSON(200, gin.H{
	// 		"message": "success",
	// 	})
	// })
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade the HTTP connection to a WebSocket connection for the client
		clientConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to WebSocket for client:", err)
			return
		}
		defer clientConn.Close()

		go proxyMessages(clientConn, serverConn) // From client to server
		proxyMessages(serverConn, clientConn)    // From server to client
	})

	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil) // listen and serve on 0.0.0.0:port
}

func proxyMessages(sender, receiver *websocket.Conn) {
	for {
		_, msg, err := sender.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		err = receiver.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}

func broadcastToClients(serverConn *websocket.Conn) {
	for {
		select {
		case msg := <-messageChannel:
			// Forward message to WebSocket server
			err := serverConn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("Error broadcasting message to WebSocket server:", err)
				break
			}
		}
	}
}
