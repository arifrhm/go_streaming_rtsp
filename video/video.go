package video

import (
	"fmt"
	"io"
	"os/exec"
)

// Constants
const (
	ffmpegPath = "ffmpeg"         // Path to FFmpeg binary
	rtspURL    = "rtsp://rtspstream:bf773dee1bab505160ae95a88c1d7585@zephyr.rtsp.stream/movie"
)

// StartFFmpegProcess starts FFmpeg as a continuous process and returns the command and reader.
func StartFFmpegProcess() (*exec.Cmd, io.Reader, error) {
	cmd := exec.Command(ffmpegPath,
		"-rtsp_transport", "tcp",
		"-i", rtspURL,
		"-f", "mjpeg",  // Output format MJPEG for continuous JPEG frames
		"-vf", "fps=20", // Adjust FPS to control frame rate
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
