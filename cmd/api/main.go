package main

import (
	"log"

	"github.com/talhag3/todoapp/internal/config"
)

func main() {
	conf, err := config.LoadConf()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Println(config.ConnectionString(conf))
}
