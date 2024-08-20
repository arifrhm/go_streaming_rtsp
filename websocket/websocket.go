package websocket

import (
	"fmt"
	"io"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/arifrhm/go_streaming_rtsp/video" // Import the video package
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// HandleWebSocketConnection handles WebSocket connections and streams frames.
func HandleWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}
	defer conn.Close()

	cmd, frameStream, err := video.StartFFmpegProcess()
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
				cmd, frameStream, err = video.StartFFmpegProcess()
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
			err = conn.WriteMessage(websocket.BinaryMessage, buffer[:n])
			if err != nil {
				fmt.Println("Error writing message:", err)
				return
			}
		}
	}
}
