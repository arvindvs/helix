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
	"github.com/arvindvs/helix/internal/stt"

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

	sttProvider, err := stt.NewProvider()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = sttProvider.StartListening(ctx)
	if err != nil {
		log.Fatal(err)
	}

	rec, err := audio.NewRecorder(16000, 8000)
	if err != nil {
		panic(err)
	}

	if err := rec.Start(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Recording... Press Ctrl+C to stop.")
	
	for {
		select {
		case <-sigs:
			slog.Info("stopping recording...")

			err = rec.Stop()
			checkErr(err)

			err = sttProvider.StopListening(ctx)
			checkErr(err)

			os.Exit(0)
		default:
			b, err := rec.Read()
			checkErr(err)

			// Send partial audio samples.
			err = sttProvider.SendAudio(ctx, b)
			checkErr(err)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		slog.Error("Something bad happened", "err", err)
		os.Exit(1)
	}
}