package speech

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"cloud.google.com/go/speech/apiv1/speechpb"
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	texttospeechpb "cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
)

type Recognizer struct {
	client *Client
}

func NewRecognizer(client *Client) *Recognizer {
	return &Recognizer{client: client}
}

func (r *Recognizer) RecognizeSpeech(ctx context.Context, audioSource <-chan []byte) error {
	stream, err := r.client.StreamingRecognize(ctx)
	if err != nil {
		return fmt.Errorf("failed to create stream: %v", err)
	}

	// Send the initial configuration request
	if err := stream.Send(&speechpb.StreamingRecognizeRequest{
		StreamingRequest: &speechpb.StreamingRecognizeRequest_StreamingConfig{
			StreamingConfig: &speechpb.StreamingRecognitionConfig{
				Config: &speechpb.RecognitionConfig{
					Encoding:        speechpb.RecognitionConfig_LINEAR16,
					SampleRateHertz: 16000,
					LanguageCode:    "en-US",
				},
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to send config: %v", err)
	}

	go r.sendAudio(stream, audioSource)
	
	return r.receiveResults(stream)
}

func (r *Recognizer) sendAudio(stream speechpb.Speech_StreamingRecognizeClient, audioSource <-chan []byte) {
	for audioData := range audioSource {
		if err := stream.Send(&speechpb.StreamingRecognizeRequest{
			StreamingRequest: &speechpb.StreamingRecognizeRequest_AudioContent{
				AudioContent: audioData,
			},
		}); err != nil {
			log.Printf("Could not send audio: %v", err)
			return
		}
	}
}

func (r *Recognizer) receiveResults(stream speechpb.Speech_StreamingRecognizeClient) error {
	var transcript strings.Builder
	var lastResultTime time.Time

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("cannot stream results: %v", err)
		}
		if err := resp.Error; err != nil {
			return fmt.Errorf("could not recognize: %v", err)
		}
		for _, result := range resp.Results {
			fmt.Println("User: ", result.Alternatives[0].Transcript)
			if result.IsFinal {
				transcript.WriteString(result.Alternatives[0].Transcript)
				transcript.WriteString(" ")
			}
		}

		// Check if there is a pause of 2 seconds or more
		if time.Since(lastResultTime) >= 2*time.Second {
			text := transcript.String()
			if text != "" {
				// Send the transcribed text as an API request
				if err := r.sendTranscriptRequest(text); err != nil {
					log.Printf("Error sending transcript request: %v", err)
				}
				transcript.Reset()
			}
			lastResultTime = time.Now()
		}
	}

	// Send any remaining transcribed text
	text := transcript.String()
	if text != "" {
		if err := r.sendTranscriptRequest(text); err != nil {
			log.Printf("Error sending transcript request: %v", err)
		}
	}

	return nil
}

func (r *Recognizer) sendTranscriptRequest(text string) error {
	url := "https://api.groq.com/openai/v1/chat/completions"

	requestBody := map[string]interface{}{
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": text,
			},
		},
		"model": "llama3-8b-8192",
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("GROQ_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		 // Parse the error response body
		var errorBody struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorBody); err != nil {
			return fmt.Errorf("request failed with status: %s", resp.Status)
		}
		return fmt.Errorf("request failed with status: %s and an error: %s", resp.Status, errorBody.Error.Message)
	}

	// Parse the JSON response
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return fmt.Errorf("failed to parse response body: %v", err)
	}

	// Extract the "content" string from the "message" object
	choices := responseBody["choices"].([]interface{})
	message := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	content := message["content"].(string)

	// Print the "content" string to the console
	fmt.Println("Assistant:", content)

	// Run speech synthesis and playback in a separate goroutine
	go func() {
		audioContent, err := synthesizeSpeech(context.Background(), content)
		if err != nil {
			log.Printf("failed to synthesize speech: %v", err)
			return
		}

		cmd := exec.Command("ffplay", "-nodisp", "-autoexit", "-")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Printf("failed to create stdin pipe: %v", err)
			return
		}
		defer stdin.Close()

		if err := cmd.Start(); err != nil {
			log.Printf("failed to start ffplay: %v", err)
			return
		}

		if _, err := io.Copy(stdin, bytes.NewReader(audioContent)); err != nil {
			log.Printf("failed to write audio data to ffplay: %v", err)
			return
		}

		if err := cmd.Wait(); err != nil {
			log.Printf("failed to wait for ffplay: %v", err)
			return
		}
	}()

	return nil
}

func synthesizeSpeech(ctx context.Context, text string) ([]byte, error) {
	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create text-to-speech client: %v", err)
	}
	defer client.Close()

	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "en-US",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_MALE,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_LINEAR16,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to synthesize speech: %v", err)
	}

	return resp.AudioContent, nil
}