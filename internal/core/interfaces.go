package core

// Define interfaces for various components
// This allows for easier testing and modularity

type AudioRecorder interface {
	// Define methods for audio recording
}

type AudioPlayer interface {
	// Define methods for audio playback
}

type WakeWordDetector interface {
	// Define methods for wake word detection
}

type NLUProcessor interface {
	// Define methods for natural language understanding
}

type TTSSynthesizer interface {
	// Define methods for text-to-speech synthesis
}

type LED interface {
	// Define methods for LED control
}