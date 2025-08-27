package server

import (
	"fmt"
	"net/http"
)

func setupRoutes(server *Server) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", server.handleHealth)
	mux.HandleFunc("/share", server.handleShare)
	mux.HandleFunc("/get/", server.handleGet)
	mux.HandleFunc("/stats", server.handleStats)

	server.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", server.config.P2PPort),
		Handler: mux,
	}
}
