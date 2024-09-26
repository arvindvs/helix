package stt

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/AssemblyAI/assemblyai-go-sdk"
)

const inactivityThreshold = 10 * time.Second // Define the inactivity threshold constant

type AssemblyAISTT struct {
	client         *assemblyai.RealTimeClient
	transcriber    *assemblyai.RealTimeTranscriber
	sampleRate     int
	lastSpeechTime time.Time
	mu             sync.Mutex
	OnTranscript   func(string)
}

func NewAssemblyAISTT(onTranscript ...func(string)) (*AssemblyAISTT, error) {
	sampleRate := 16000

	var transcriptFunc func(string)
	if len(onTranscript) > 0 {
		transcriptFunc = onTranscript[0]
	}

	stt := &AssemblyAISTT{
		sampleRate:   sampleRate,
		OnTranscript: transcriptFunc,
	}

	transcriber := &assemblyai.RealTimeTranscriber{
		OnSessionBegins: func(event assemblyai.SessionBegins) { },
		OnSessionTerminated: func(event assemblyai.SessionTerminated) { },
		OnFinalTranscript: func(transcript assemblyai.FinalTranscript) {
			fmt.Println(transcript.Text)
			stt.updateLastSpeechTime()
			if stt.OnTranscript != nil {
				stt.OnTranscript(transcript.Text)
			}
		},
		OnPartialTranscript: func(transcript assemblyai.PartialTranscript) {
			fmt.Printf("%s\r", transcript.Text)
			stt.updateLastSpeechTime()
		},
		OnError: func(err error) {
			slog.Error("System failure:", "error", err.Error())
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

	stt.client = client
	stt.transcriber = transcriber

	return stt, nil
}

func (a *AssemblyAISTT) StartListening(ctx context.Context) error {
	a.updateLastSpeechTime()
	return a.client.Connect(ctx)
}

func (a *AssemblyAISTT) StopListening(ctx context.Context) error {
	return a.client.Disconnect(ctx, true)
}

func (a *AssemblyAISTT) SendAudio(ctx context.Context, audio []byte) error {
	return a.client.Send(ctx, audio)
}

func (a *AssemblyAISTT) IsSpeechDetected() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return time.Since(a.lastSpeechTime) < inactivityThreshold // Use the constant here
}

func (a *AssemblyAISTT) updateLastSpeechTime() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.lastSpeechTime = time.Now()
}
