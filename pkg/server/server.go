package server

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type PublishMessage struct {
	Data string `json:"data"`
}

var messageChannel = make(chan []byte) // Channel to send data to the WebSocket server

func NewWebSocketProxy(wsUrl *url.URL) http.HandlerFunc {
	// Create a reverse proxy director
	director := func(req *http.Request) {
		req.URL.Scheme = wsUrl.Scheme
		req.URL.Host = wsUrl.Host
		req.URL.Path = wsUrl.Path
	}

	// Create the reverse proxy
	proxy := &httputil.ReverseProxy{Director: director, Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}}

	return func(w http.ResponseWriter, r *http.Request) {
		// Proxy the WebSocket connection
		proxy.ServeHTTP(w, r)
	}
}

func writeMessageToWebsocket(serverConn *websocket.Conn) {
	for {
		msg, ok := <-messageChannel
		if ok {
			// Forward message to WebSocket server
			err := serverConn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("Error broadcasting message to WebSocket server:", err)
				break
			}
		}
	}
}

// Start my server instance
func StartServer(port, wsServerUrl string) {
	log.Printf("Client listen and serve at %s\n", port)

	// Parse the internal WebSocket URL
	wsUrl, err := url.Parse(wsServerUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Establish a WebSocket connection to the WebSocket server
	// to be use for publishing data to the websocket server.
	websocketConn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/ws", wsUrl.Host), nil)
	if err != nil {
		log.Fatal("WebSocket dial:", err)
	}
	defer websocketConn.Close()

	// write message to the main websocket server in order to broadcast message to all other server instances.
	go writeMessageToWebsocket(websocketConn)

	wsProxy := NewWebSocketProxy(wsUrl)

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World")
	})

	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		data := PublishMessage{Data: ""}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Invalid post data", http.StatusBadRequest)
			return
		}

		fmt.Println("Writing data: ", data.Data)

		messageChannel <- []byte(data.Data) // forwading message to websocket server

		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "OK",
			"message": "Data successfully published",
		})
	})

	http.HandleFunc("/ws", wsProxy)

	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil) // listen and serve on 0.0.0.0:port
}
