package main

import (
	"net/http"

	"github.com/vlad-sidius/go-url-shortener/internal/handler"
)

func main() {
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
