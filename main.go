package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"io"
	"github.com/gorilla/websocket"
)

// Constants
const (
	ffmpegPath = "ffmpeg"         // Path to FFmpeg binary
	rtspURL    = "rtsp://rtspstream:bf773dee1bab505160ae95a88c1d7585@zephyr.rtsp.stream/movie"
	wsPort      = ":8330"          // Port for WebSocket server
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// StartFFmpegProcess starts FFmpeg as a continuous process and returns the command and reader.
func startFFmpegProcess() (*exec.Cmd, io.Reader, error) {
	cmd := exec.Command(ffmpegPath,
		"-rtsp_transport", "tcp",
		"-i", rtspURL,
		"-f", "mjpeg",  // Output format MJPEG for continuous JPEG frames
		"-vf", "fps=10", // Adjust FPS to control frame rate
		"-")

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("Error creating stdout pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("Error starting FFmpeg: %v", err)
	}

	return cmd, stdoutPipe, nil
}

// HandleWebSocketConnection handles WebSocket connections and streams frames.
func handleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}
	defer conn.Close()

	cmd, frameStream, err := startFFmpegProcess()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cmd.Process.Kill()

	buffer := make([]byte, 1024*1024) // Buffer size for MJPEG frames

	for {
		n, err := frameStream.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF reached. Restarting FFmpeg process...")
				cmd.Process.Kill()
				cmd, frameStream, err = startFFmpegProcess()
				if err != nil {
					fmt.Println(err)
					return
				}
				continue
			}
			fmt.Println("Error reading frame data:", err)
			return
		}

		if n > 0 {
			// Send the MJPEG frame over WebSocket
			err = conn.WriteMessage(websocket.BinaryMessage, buffer[:n])
			if err != nil {
				fmt.Println("Error writing message:", err)
				return
			}
		}
	}
}

// IndexHandler serves the index.html file.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// StartHTTPServer starts the HTTP server for serving the HTML file and WebSocket connections.
func startHTTPServer() {
	http.HandleFunc("/", indexHandler) // Serve index.html
	http.HandleFunc("/ws", handleWebSocketConnection) // WebSocket for streaming

	fmt.Printf("HTTP server started at http://localhost%s\n", wsPort)
	if err := http.ListenAndServe(wsPort, nil); err != nil {
		fmt.Printf("Error starting HTTP server: %v\n", err)
	}
}

func main() {
	startHTTPServer()
}
