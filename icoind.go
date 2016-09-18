package main

import (
	"log"
	"os"
)

var (
	conf *config
)

func main() {
	conf, err := loadConfig()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	server := newServer(conf)

	server.Start()
}
