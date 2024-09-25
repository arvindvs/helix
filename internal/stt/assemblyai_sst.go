package stt

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/AssemblyAI/assemblyai-go-sdk"
)

type AssemblyAISTT struct {
	client      *assemblyai.RealTimeClient
	transcriber *assemblyai.RealTimeTranscriber
	sampleRate  int
}

func NewAssemblyAISTT() (*AssemblyAISTT, error) {
	sampleRate := 16000

	transcriber := &assemblyai.RealTimeTranscriber{
		OnSessionBegins: func(event assemblyai.SessionBegins) {
			slog.Info("session begins")
		},
		OnSessionTerminated: func(event assemblyai.SessionTerminated) {
			slog.Info("session terminated")
		},
		OnFinalTranscript: func(transcript assemblyai.FinalTranscript) {
			fmt.Println(transcript.Text)
		},
		OnPartialTranscript: func(transcript assemblyai.PartialTranscript) {
			fmt.Printf("%s\r", transcript.Text)
		},
		OnError: func(err error) {
			slog.Error("Something bad happened", "err", err)
		},
	}

	apiKey := os.Getenv("ASSEMBLYAI_API_KEY")
	
	if apiKey == "" {
		return nil, fmt.Errorf("ASSEMBLYAI_API_KEY is not set")
	}

	client := assemblyai.NewRealTimeClientWithOptions(
		assemblyai.WithRealTimeAPIKey(apiKey),
		assemblyai.WithRealTimeSampleRate(sampleRate),
		assemblyai.WithRealTimeTranscriber(transcriber),
	)

	return &AssemblyAISTT{
		client:      client,
		transcriber: transcriber,
		sampleRate:  sampleRate,
	}, nil
}

func (a *AssemblyAISTT) StartListening(ctx context.Context) error {
	return a.client.Connect(ctx)
}

func (a *AssemblyAISTT) StopListening(ctx context.Context) error {
	return a.client.Disconnect(ctx, true)
}

func (a *AssemblyAISTT) SendAudio(ctx context.Context, audio []byte) error {
	return a.client.Send(ctx, audio)
}
