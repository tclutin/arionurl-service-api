package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address string        `yaml:"address" env-default:"localhost:8080"`
	Timeout time.Duration `yaml:"timeout" env-default:"2s"`
}

type DBConfig struct {
}

func MustLoad() *Config {
	var config Config

	path := os.Getenv("ARIONURL_CONFIG")
	if path == "" {
		log.Fatalln("specify the correct path to the config")
	}

	if _, err := os.ReadFile(path); err != nil {
		log.Fatalln(err)
	}

	if err := cleanenv.ReadConfig(path, &config); err != nil {
		log.Fatalln(err)
	}

	return &config
}
