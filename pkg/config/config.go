package config

import (
	"github.com/BurntSushi/toml"
)

// Config represents mix of settings for the app.
type Config struct {
	Log      Log      `toml:"log"`
	Server   Server   `toml:"server"`
	OpenALPR OpenALPR `toml:"openalpr"`
}

// Server represents settings for creating http server.
type Server struct {
	ShutdownTimeout Duration `toml:"shutdown_timeout"`
	ReadTimeout     Duration `toml:"read_timeout"`
	WriteTimeout    Duration `toml:"write_timeout"`
	IdleTimeout     Duration `toml:"idle_timeout"`
}

type Log struct {
	Level string `toml:"level"`
	Mode  string `toml:"mode"`
}

type OpenALPR struct {
	Country    string `toml:"country"`
	ConfigFile string `toml:"config_file"`
	RuntimeDir string `toml:"runtime_dir"`
	MaxNumber  int    `toml:"max_number"`
}

// New creates reads application configuration from the file.
func New(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
