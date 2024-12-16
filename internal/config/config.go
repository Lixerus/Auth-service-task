package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type TokenSecret struct {
	Text []byte
}

func (t *TokenSecret) UnmarshalText(text []byte) error {
	t.Text = text
	return nil
}

type Config struct {
	HOST          string      `env:"PG_HOST"`
	PORT          string      `env:"PG_PORT"`
	USERNAME      string      `env:"PG_USERNAME"`
	PASSWORD      string      `env:"PG_PASSWORD"`
	NAME          string      `env:"PG_DATABASE"`
	ACCESS_SECRET TokenSecret `env:"ACCESS_TOKEN_SECRET"`
}

var DBConfig *Config

func LoadEnv() {
	err := godotenv.Load("../.././internal/config/db.env")
	if err != nil {

		log.Fatal(err.Error())
	}
}

func LoadConfig() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal("Unable to load config" + err.Error())
	}
	DBConfig = &cfg
}
