package config

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Mode     string `env:"MODE"`
	HttpPort string `env:"HTTP_PORT"`

	ImmudbUser     string `env:"IMMUDB_USER"`
	ImmudbPwd      string `env:"IMMUDB_PWD"`
	ImmudbDatabase string `env:"IMMUDB_DB_NAME"`
	ImmudbHost     string `env:"IMMUDB_HOST"`
	ImmudbPort     int    `env:"IMMUDB_PORT"`
	ImmudbSsl      bool   `env:"IMMUDB_SSL"`
	DbTable        string `env:"DB_TABLE"`

	MaxLifetimeInMinutes int `env:"MAX_LIFETIME_IN_MINUTES"`
	MaxConnections       int `env:"MAX_CONNECTIONS"`
	ConnIdle             int `env:"CONN_IDLE"`
}

func LoadEnvConfig() Config {
	var config Config

	godotenv.Load(".env")
	ctx := context.Background()
	if err := envconfig.Process(ctx, &config); err != nil {
		panic(err)
	}

	return config
}
