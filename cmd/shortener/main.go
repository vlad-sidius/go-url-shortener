package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vlad-sidius/go-url-shortener/internal/config"
	"github.com/vlad-sidius/go-url-shortener/internal/handler"
	mw "github.com/vlad-sidius/go-url-shortener/internal/handler/middleware"
	"github.com/vlad-sidius/go-url-shortener/internal/logger"
	"github.com/vlad-sidius/go-url-shortener/internal/repository"
	"github.com/vlad-sidius/go-url-shortener/internal/service"
)

func main() {
	conf := config.InitConfig()

	zLogger, err := logger.NewLogger()
	if err != nil {
		log.Fatal("Failed to initialize logger")
	}

	memRepo := repository.NewMemURLRepo(zLogger)
	if err := memRepo.Load(conf.URLServiceConf.StoragePath); err != nil {
		zLogger.Error("Failed to load data from file. Continue with empty storage.")
	}

	hashGen := service.NewRandomSlugGenerator()
	urlService := service.NewURLServiceLive(zLogger, &conf.URLServiceConf, memRepo, hashGen)
	urlHandler := handler.NewURLHandler(zLogger, urlService)

	router := gin.New()
	router.Use(mw.LoggerMiddleware(zLogger), mw.GzipMiddleware(zLogger), gin.Recovery())
	urlHandler.RegisterRoutes(router)

	err = router.Run(conf.ServerConf.Address)
	if err != nil {
		log.Fatalf("Failed to start server %v\n", err)
	}
}
