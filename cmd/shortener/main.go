package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/vlad-sidius/go-url-shortener/internal/config"
	"github.com/vlad-sidius/go-url-shortener/internal/handler"
	mw "github.com/vlad-sidius/go-url-shortener/internal/handler/middleware"
	"github.com/vlad-sidius/go-url-shortener/internal/logger"
	"github.com/vlad-sidius/go-url-shortener/internal/repository"
	"github.com/vlad-sidius/go-url-shortener/internal/service"
	"go.uber.org/zap"
)

func main() {
	conf := config.InitConfig()

	zLogger, err := logger.NewLogger()
	if err != nil {
		log.Fatal("Failed to initialize logger")
	}

	fileRepo := repository.NewFileURLRepo(zLogger, conf.URLServiceConf.StoragePath)
	hashGen := service.NewRandomSlugGenerator()
	urlService := service.NewURLServiceLive(zLogger, &conf.URLServiceConf, fileRepo, hashGen)
	urlHandler := handler.NewURLHandler(zLogger, urlService)

	router := gin.New()
	router.Use(mw.LoggerMiddleware(zLogger), mw.GzipMiddleware(zLogger), gin.Recovery())
	urlHandler.RegisterRoutes(router)

	err = router.Run(conf.ServerConf.Address)
	if err != nil {
		zLogger.Error("Failed to start server", zap.Error(err))
		os.Exit(1)
	}
}
