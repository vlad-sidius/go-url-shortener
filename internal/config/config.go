package config

import "flag"

type Config struct {
	Address string
	BaseURL string
}

func ParseCliArgs() *Config {
	var config Config

	flag.StringVar(&config.Address, "a", "localhost:8080", "Server address")
	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "Base part for short URL")

	flag.Parse()

	return &config
}
