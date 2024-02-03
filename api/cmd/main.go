package main

import (
	"fmt"
	"github.com/tclutin/ArionURL/internal/config"
)

func main() {
	//Initializing the config
	cfg := config.MustLoad()
	fmt.Println(cfg)
}
