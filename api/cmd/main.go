package main

import (
	"github.com/tclutin/ArionURL/internal/config"
	"github.com/tclutin/ArionURL/pkg/logging"
)

func main() {
	//Initializing the config
	cfg := config.MustLoad()

	//Initializing the logger
	logger := logging.InitSlog(cfg.Env)
}
