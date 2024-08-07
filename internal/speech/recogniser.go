package speech

import (
	"context"
	"fmt"
	"io"
	"log"

	"cloud.google.com/go/speech/apiv1/speechpb"
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
			fmt.Printf("Result: %+v\n", result)
		}
	}
	return nil
}