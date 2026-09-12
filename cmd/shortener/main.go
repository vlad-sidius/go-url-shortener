package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vlad-sidius/go-url-shortener/internal/config"
	"github.com/vlad-sidius/go-url-shortener/internal/handler"
	"github.com/vlad-sidius/go-url-shortener/internal/logger"
	"github.com/vlad-sidius/go-url-shortener/internal/repository"
	"github.com/vlad-sidius/go-url-shortener/internal/repository/middleware"
	"github.com/vlad-sidius/go-url-shortener/internal/service"
)

func main() {
	conf := config.InitConfig()

	zLogger, err := logger.NewLogger()
	if err != nil {
		log.Fatal("Failed to initialize logger")
	}

	memRepo := repository.NewMemURLRepo()
	hashGen := service.NewRandomSlugGenerator()
	urlService := service.NewURLServiceLive(&conf.URLServiceConf, memRepo, hashGen)
	urlHandler := handler.NewURLHandler(zLogger, urlService)

	router := gin.New()
	router.Use(middleware.LoggerMiddleware(zLogger), middleware.GzipMiddleware(zLogger), gin.Recovery())
	urlHandler.RegisterRoutes(router)

	err = router.Run(conf.ServerConf.Address)
	if err != nil {
		log.Fatalf("Failed to start server %v\n", err)
	}
}
