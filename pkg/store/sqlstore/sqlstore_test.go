package sqlstore_test

import (
	"os"
	"testing"

	"github.com/opencars/alpr/pkg/config"
)

var (
	conf *config.Database
)

func TestMain(m *testing.M) {
	conf = &config.Database{
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     5432,
		User:     "postgres",
		Password: "password",
		Name:     "alpr",
		SSLMode:  "disable",
	}

	if conf.Host == "" {
		conf.Host = "127.0.0.1"
	}

	code := m.Run()
	os.Exit(code)
}
