package stt

import (
	"context"
)

type STT interface {
	StartListening(ctx context.Context) error
	StopListening(ctx context.Context) error
	SendAudio(ctx context.Context, audio []byte) error
}

func NewProvider() (STT, error) {
	return NewAssemblyAISTT()
}
