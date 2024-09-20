package configs

// Import necessary packages

type Config struct {
	Audio struct {
		InputDevice  string `json:"input_device"`
		OutputDevice string `json:"output_device"`
	} `json:"audio"`
	WakeWord struct {
		ModelPath string `json:"model_path"`
	} `json:"wake_word"`
	NLU struct {
		ModelPath string `json:"model_path"`
	} `json:"nlu"`
	TTS struct {
		Voice string `json:"voice"`
	} `json:"tts"`
	APIKeys struct {
		WeatherAPI string `json:"weather_api"`
		NewsAPI    string `json:"news_api"`
	} `json:"api_keys"`
	Hardware struct {
		LEDPin int `json:"led_pin,omitempty"`
	} `json:"hardware,omitempty"`
}

func LoadConfig(filename string) (*Config, error) {
	// Load and parse configuration from file
}