package api

import (
	"net/http"
	"time"

	"github.com/m1kx/go-vtr-backend/pkg/utils/health"
)

func RunServer() {
	http.HandleFunc("/health", health.Health)
	server := &http.Server{
		Addr:         ":9999",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	go server.ListenAndServe()
}
