package config

import (
	"flag"
	"os"
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
	updateConfigWithEnvOverrides(config)

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

func updateConfigWithEnvOverrides(conf *Config) {
	address, ok := os.LookupEnv(addressEnv)
	if ok {
		conf.ServerConf.Address = address
	}

	baseURL, ok := os.LookupEnv(baseURLEnv)
	if ok {
		conf.URLServiceConf.BaseURL = baseURL
	}

	storagePath, ok := os.LookupEnv(storagePathEnv)
	if ok {
		conf.URLServiceConf.StoragePath = storagePath
	}
}
