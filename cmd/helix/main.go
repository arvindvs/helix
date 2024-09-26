package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/arvindvs/helix/internal/audio"
	"github.com/arvindvs/helix/internal/core"
	"github.com/arvindvs/helix/internal/stt"
	"github.com/arvindvs/helix/internal/wakeword"

	"github.com/gordonklaus/portaudio"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("🧬 Starting Helix")
	fmt.Println("🔍 Loading environment variables...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) 
	
	err = portaudio.Initialize()
	if err != nil {
		panic(err)
	}
	defer portaudio.Terminate()

	wakeWordDetector, err := wakeword.NewDetector("path/to/porcupine_model.pv", "path/to/helix_keyword.ppn", 0.5)
	if err != nil {
		log.Fatal(err)
	}
	defer wakeWordDetector.Close()

	assistant, err := core.NewAssistant()
	if err != nil {
		log.Fatal(err)
	}

	sttProvider, err := stt.NewAssemblyAISTT(func(text string) {
		assistant.ProcessText(context.Background(), text)
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()


	rec, err := audio.NewRecorder(16000, 8192)
	if err != nil {
		panic(err)
	}

	if err := rec.Start(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening for wake word 'Helix'... Press Ctrl+C to stop.")
	
	isListening := false

	for {
		select {
		case <-sigs:
			err = rec.Stop()
			checkErr(err)

			err = sttProvider.StopListening(ctx)
			checkErr(err)

			os.Exit(0)
		default:
			b, err := rec.Read()
			checkErr(err)

			if !isListening {
				for i := 0; i < len(b); i += 1024 {
					end := i + 1024
					if end > len(b) {
						end = len(b)
					}
					chunk := b[i:end]
					detected, err := wakeWordDetector.Process(chunk)
					checkErr(err)

					if detected {
						slog.Info("Wake word detected. Listening...")
						fmt.Println("Your input:")
						err = sttProvider.StartListening(ctx)
						checkErr(err)
						isListening = true
						break
					}
				}
			} else {
				// Send partial audio samples to STT provider
				err = sttProvider.SendAudio(ctx, b)
				checkErr(err)

				if !sttProvider.IsSpeechDetected() {
					slog.Info("No speech detected for 10 seconds. Now listening for wake word...")
					err = sttProvider.StopListening(ctx)
					checkErr(err)
					isListening = false
					rec.Stop()
					rec.Start()
				}
			}
		}
	}
}

func checkErr(err error) {
	if err != nil {
		slog.Error("System failure", "err", err)
		os.Exit(1)
	}
}