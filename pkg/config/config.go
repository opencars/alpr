package config

import (
	"github.com/BurntSushi/toml"
)

// Config represents mix of settings for the app.
type Config struct {
	Log      Log      `toml:"log"`
	Server   Server   `toml:"server"`
	OpenALPR OpenALPR `toml:"openalpr"`
	S3       S3       `toml:"s3"`
}

// Server represents settings for creating http server.
type Server struct {
	ShutdownTimeout Duration `toml:"shutdown_timeout"`
	ReadTimeout     Duration `toml:"read_timeout"`
	WriteTimeout    Duration `toml:"write_timeout"`
	IdleTimeout     Duration `toml:"idle_timeout"`
}

// Log represents settings for application logger.
type Log struct {
	Level string `toml:"level"`
	Mode  string `toml:"mode"`
}

// OpenALPR represents settings for openalpr car plates recognizer.
type OpenALPR struct {
	Country    string `toml:"country"`
	ConfigFile string `toml:"config_file"`
	RuntimeDir string `toml:"runtime_dir"`
	MaxNumber  int    `toml:"max_number"`
}

// S3 represents settings for s3 storage.
type S3 struct {
	Endpoint        string `toml:"endpoint"`
	AccessKeyID     string `toml:"access_key_id"`
	SecretAccessKey string `toml:"secret_access_key"`
	SSL             bool   `toml:"ssl"`
	Bucket          string `toml:"bucket"`
}

// New creates reads application configuration from the file.
func New(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
