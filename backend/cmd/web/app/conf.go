package app

type conf struct {
	Port          string `yaml:"port"`
	DSN           string `yaml:"dsn"`
	TimezoneLoc   string `yaml:"timezone_loc"`
	SessionSecret string `yaml:"session_secret"`
	Storage_path  string `yaml:"storage_path"`
	GroupID       string `yaml:"group_id"`
	KafkaHost     string `yaml:"kafka_host"`
	Topic         string `yaml:"topic"`
	STT           struct {
		PushAudioURL string `yaml:"push_audio_url"`
		DataAudioURl string `yaml:"data_audio_url"`
		Token        string `yaml:"token"`
		TryInterval  int    `yaml:"try_interval"`
	} `yaml:"stt"`
	TTS struct {
		URL   string `yaml:"url"`
		Token string `yaml:"token"`
	} `yaml:"tts"`
	GPT struct {
		Token string `yaml:"token"`
		Url   string `yaml:"url"`
	} `yaml:"gpt"`
	OAuth2 struct {
		RedirectURL      string   `yaml:"redirect_url"`
		ClientID         string   `yaml:"client_id"`
		ClientSecret     string   `yaml:"client_secret"`
		Scopes           []string `yaml:"scopes"`
		UserInfoEndpoint string   `yaml:"user_info_endpoint"`
	} `yaml:"google"`
}
