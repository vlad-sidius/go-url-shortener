package config

import "flag"

type ServerConfig struct {
	Address string
}

type URLServiceConfig struct {
	BaseURL string
}

type Config struct {
	ServerConf     ServerConfig
	URLServiceConf URLServiceConfig
}

func NewConfig(address, baseURL string) *Config {
	return &Config{
		ServerConf:     ServerConfig{address},
		URLServiceConf: URLServiceConfig{baseURL},
	}
}

func ParseCliArgs() *Config {
	var serverConf ServerConfig
	var urlServiceConf URLServiceConfig

	flag.StringVar(&serverConf.Address, "a", "localhost:8080", "Server address")
	flag.StringVar(&urlServiceConf.BaseURL, "b", "http://localhost:8080", "Base part for short URL")

	flag.Parse()

	return &Config{serverConf, urlServiceConf}
}
