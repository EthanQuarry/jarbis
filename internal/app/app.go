package app

import (
	"context"
	"time"

	"github.com/EthanQuarry/jarbis/internal/speech"
)

type App struct {
	speechClient *speech.Client
	recognizer   *speech.Recognizer
	recorder     *speech.Recorder
}

func NewApp(ctx context.Context) (*App, error) {
	client, err := speech.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &App{
		speechClient: client,
		recognizer:   speech.NewRecognizer(client),
		recorder:     speech.NewRecorder(100 * time.Millisecond),
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	defer a.speechClient.Close()

	audioCh := make(chan []byte)
	go a.recorder.Record(audioCh)

	return a.recognizer.RecognizeSpeech(ctx, audioCh)
}