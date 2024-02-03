package main

import (
	"context"
	"fmt"
	"github.com/tclutin/ArionURL/internal/config"
	"github.com/tclutin/ArionURL/pkg/client/postgresql"
	"github.com/tclutin/ArionURL/pkg/logging"
	"os"
)

func main() {
	//Initializing the config
	cfg := config.MustLoad()

	//Initializing the logger
	logger := logging.InitSlog(cfg.Env)

	//Initial the pgxpool
	pgxPool := postgresql.NewClient(context.Background(), os.Getenv("ARIONURL_DB"))
	fmt.Println(logger, pgxPool)
}
