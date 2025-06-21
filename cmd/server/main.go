package main

import (
	"log"
	"net/http"

	// internal layers
	"github.com/codex/go-web-console/internal/interface/controller"
	"github.com/codex/go-web-console/internal/interface/gateway"
	"github.com/codex/go-web-console/internal/usecase/log"
)

// main wires dependencies manually and starts the HTTP server.
func main() {
	repo := gateway.NewFileLogStore("logs.json")
	uc := log.NewViewerUsecase(repo)
	handler := controller.NewLogHandler(uc)

	mux := http.NewServeMux()
	mux.Handle("/logs/{level}", handler)

	log.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
