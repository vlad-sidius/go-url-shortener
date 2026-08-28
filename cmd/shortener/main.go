package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/vlad-sidius/go-url-shortener/internal/config"
	"github.com/vlad-sidius/go-url-shortener/internal/handler"
	"github.com/vlad-sidius/go-url-shortener/internal/repository"
)

func main() {
	conf := config.ParseCliArgs()
	memRepo := repository.NewMemURLRepo()
	urlHandler := handler.NewURLHandler(conf, memRepo)

	router := gin.Default()
	urlHandler.RegisterRoutes(router)

	err := router.Run(conf.Address)
	if err != nil {
		fmt.Printf("Failed to start server %v\n", err)
		os.Exit(1)
	}
}
