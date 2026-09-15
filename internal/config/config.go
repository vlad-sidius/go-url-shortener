package config

import (
	"flag"
	"os"
	"strings"
)

const (
	addressEnv     = "SERVER_ADDRESS"
	baseURLEnv     = "BASE_URL"
	storagePathEnv = "FILE_STORAGE_PATH"
)

type ServerConfig struct {
	Address string
}

type URLServiceConfig struct {
	BaseURL     string
	StoragePath string
}

type Config struct {
	ServerConf     ServerConfig
	URLServiceConf URLServiceConfig
}

func InitConfig() *Config {
	config := readCliArgs()
	address, baseURL, storagePath := readEnvVars()

	if len(strings.TrimSpace(address)) > 0 {
		config.ServerConf.Address = address
	}

	if len(strings.TrimSpace(storagePath)) > 0 {
		config.URLServiceConf.StoragePath = storagePath
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
	flag.StringVar(&urlServiceConf.StoragePath, "f", "fileDb.json", "Storage file path")
	flag.StringVar(&urlServiceConf.BaseURL, "b", "http://localhost:8080", "Base part for short URL")

	flag.Parse()

	return &Config{serverConf, urlServiceConf}
}

func readEnvVars() (string, string, string) {
	address := os.Getenv(addressEnv)
	baseURL := os.Getenv(baseURLEnv)
	storagePath := os.Getenv(storagePathEnv)

	return address, baseURL, storagePath
}
