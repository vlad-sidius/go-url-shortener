package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vlad-sidius/go-url-shortener/internal/handler"
	"github.com/vlad-sidius/go-url-shortener/internal/repository"
)

func main() {
	memRepo := repository.NewMemURLRepo()
	urlHandler := handler.NewURLHandler(memRepo)

	//urlHandler.RegisterRoutes(mux)
	//
	//err := http.ListenAndServe(`:8080`, mux)
	//if err != nil {
	//	panic(err)
	//}

	router := gin.Default()
	urlHandler.RegisterRoutes(router)

	_ = router.Run(`:8080`)
}
