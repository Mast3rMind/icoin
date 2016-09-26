package main

import (
	"errors"
	"flag"
)

type config struct {
	net     string
	connect string
	params  *params
}

func loadConfig() (*config, error) {
	var conf config

	flag.StringVar(&conf.net, "net", "main", "which network to connect to")
	flag.StringVar(&conf.connect, "connect", "", "which node to connect to")

	flag.Parse()

	switch conf.net {
	case "main":
		conf.params = &mainNetParams
	case "test":
		conf.params = &testNetParams
	default:
		return nil, errors.New("invalid network type")
	}

	return &conf, nil
}
