package config

import "flag"

type Config struct {
	address string
	baseURL string
}

func NewConfig(address, baseURL string) *Config {
	return &Config{address, baseURL}
}

func ParseCliArgs() *Config {
	var config Config

	flag.StringVar(&config.address, "a", "localhost:8080", "Server address")
	flag.StringVar(&config.baseURL, "b", "http://localhost:8080", "Base part for short URL")

	flag.Parse()

	return &config
}

func (c *Config) Address() string {
	return c.address
}

func (c *Config) BaseURL() string {
	return c.baseURL
}
