package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vlad-sidius/go-url-shortener/internal/config"
	"github.com/vlad-sidius/go-url-shortener/internal/handler"
	"github.com/vlad-sidius/go-url-shortener/internal/repository"
	"github.com/vlad-sidius/go-url-shortener/internal/service"
)

func main() {
	conf := config.ParseCliArgs()
	memRepo := repository.NewMemURLRepo()
	hashGen := service.NewRandomSlugGenerator()
	urlService := service.NewURLServiceLive(&conf.URLServiceConf, memRepo, hashGen)
	urlHandler := handler.NewURLHandler(urlService)

	router := gin.Default()
	urlHandler.RegisterRoutes(router)

	err := router.Run(conf.ServerConf.Address)
	if err != nil {
		log.Fatalf("Failed to start server %v\n", err)
	}
}
