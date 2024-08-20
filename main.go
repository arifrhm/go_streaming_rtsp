package main

import (
	"fmt"
	"net/http"
	"github.com/arifrhm/go_streaming_rtsp/websocket" // Import the websocket package
)

// Constants
const wsPort = ":8330" // Port for WebSocket server

// IndexHandler serves the index.html file.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// StartHTTPServer starts the HTTP server for serving the HTML file and WebSocket connections.
func startHTTPServer() {
	http.HandleFunc("/", indexHandler) // Serve index.html
	http.HandleFunc("/ws", websocket.HandleWebSocketConnection) // WebSocket for streaming

	fmt.Printf("HTTP server started at http://localhost%s\n", wsPort)
	if err := http.ListenAndServe(wsPort, nil); err != nil {
		fmt.Printf("Error starting HTTP server: %v\n", err)
	}
}

func main() {
	startHTTPServer()
}
