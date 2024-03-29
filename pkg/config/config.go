package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config represents mix of settings for the app.
type Config struct {
	Log      Log      `yaml:"log"`
	Server   Server   `yaml:"server"`
	OpenALPR OpenALPR `yaml:"openalpr"`
	S3       S3       `yaml:"s3"`
	DB       Database `yaml:"database"`
	NATS     NATS     `yaml:"nats"`
}

// Server represents settings for creating http server.
type Server struct {
	ShutdownTimeout Duration `yaml:"shutdown_timeout"`
	ReadTimeout     Duration `yaml:"read_timeout"`
	WriteTimeout    Duration `yaml:"write_timeout"`
	IdleTimeout     Duration `yaml:"idle_timeout"`
}

// Log represents settings for application logger.
type Log struct {
	Level string `yaml:"level"`
	Mode  string `yaml:"mode"`
}

// OpenALPR represents settings for openalpr car plates recognizer.
type OpenALPR struct {
	Pool       int    `yaml:"pool"`
	Country    string `yaml:"country"`
	ConfigFile string `yaml:"config_file"`
	RuntimeDir string `yaml:"runtime_dir"`
	MaxNumber  int    `yaml:"max_number"`
}

// S3 represents settings for s3 storage.
type S3 struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	SSL             bool   `yaml:"ssl"`
	Bucket          string `yaml:"bucket"`
	URL             string `yaml:"base_url"`
}

// Database contains configuration details for database.
type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"database"`
	SSLMode  string `yaml:"ssl_mode"`
}

// Address return API address in "host:port" format.
func (db *Database) Address() string {
	return db.Host + ":" + strconv.Itoa(db.Port)
}

type NodeNATS struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (node *NodeNATS) Address(user, password string) string {
	if user != "" && password != "" {
		return fmt.Sprintf("nats://%s:%s@%s:%d", user, password, node.Host, node.Port)
	}

	return fmt.Sprintf("nats://%s:%d", node.Host, node.Port)
}

// NATS contains configuration details for application event API.
type NATS struct {
	Enabled  bool       `yaml:"enabled"`
	Nodes    []NodeNATS `yaml:"nodes"`
	User     string     `yaml:"user"`
	Password string     `yaml:"password"`
}

// Address returns calculated address for connecting to NATS.
func (nats *NATS) Address() string {
	addrs := make([]string, 0)

	for _, node := range nats.Nodes {
		addrs = append(addrs, node.Address(nats.User, nats.Password))
	}

	return strings.Join(addrs, ",")
}

// New reads application configuration from specified file path.
func New(path string) (*Config, error) {
	var config Config

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
