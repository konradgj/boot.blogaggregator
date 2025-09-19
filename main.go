package main

import (
	"fmt"
	"log"
	"os"

	"github.com/konradgj/boot.blogaggregator/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = cfg.SetUser("lane")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	cfg, err = config.LoadConfig()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println(cfg)
}
