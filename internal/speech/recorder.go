package speech

import (
	"fmt"
	"os/exec"
	"time"
)

type Recorder struct {
	duration time.Duration
}

func NewRecorder(duration time.Duration) *Recorder {
	return &Recorder{duration: duration}
}

func (r *Recorder) Record(ch chan<- []byte) {
	fmt.Println("Recording started. Press Ctrl+C to stop.")
	for {
		cmd := exec.Command("ffmpeg",
			"-f", "dshow",
			"-i", "audio=Microphone (3- USB PnP Audio Device)",
			 "-acodec", "pcm_s16le",
			"-ar", "16000",
			"-ac", "1",
			"-f", "s16le",
			"pipe:1")

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Printf("Error creating stdout pipe: %v\n", err)
			continue
		}

		if err := cmd.Start(); err != nil {
			fmt.Printf("Error starting ffmpeg: %v\n", err)
			continue
		}

		buffer := make([]byte, 1024)
		for {
			n, err := stdout.Read(buffer)
			if err != nil {
				break
			}
			ch <- buffer[:n]
		}

		if err := cmd.Wait(); err != nil {
			// Check if the error is due to interruption (Ctrl+C)
			if err.Error() == "signal: interrupt" {
				fmt.Println("Recording stopped.")
				return
			}
			fmt.Printf("Error waiting for ffmpeg to finish: %v\n", err)
		}
	}
}