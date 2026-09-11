package config

import (
	"flag"
	"os"
	"strings"
)

const (
	addressEnv = "SERVER_ADDRESS"
	baseURLEnv = "BASE_URL"
)

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

func InitConfig() *Config {
	config := readCliArgs()
	address, baseURL := readEnvVars()

	if len(strings.TrimSpace(address)) > 0 {
		config.ServerConf.Address = address
	}

	if len(strings.TrimSpace(baseURL)) > 0 {
		config.URLServiceConf.BaseURL = baseURL
	}

	return config
}

func readCliArgs() *Config {
	var serverConf ServerConfig
	var urlServiceConf URLServiceConfig

	flag.StringVar(&serverConf.Address, "a", "localhost:8080", "Server address")
	flag.StringVar(&urlServiceConf.BaseURL, "b", "http://localhost:8080", "Base part for short URL")

	flag.Parse()

	return &Config{serverConf, urlServiceConf}
}

func readEnvVars() (string, string) {
	address := os.Getenv(addressEnv)
	baseURL := os.Getenv(baseURLEnv)

	return address, baseURL
}
