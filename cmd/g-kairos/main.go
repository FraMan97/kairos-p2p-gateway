package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/FraMan97/kairos-p2p-gateway/internal/api"
	"github.com/FraMan97/kairos-p2p-gateway/internal/config"
)

func main() {
	config.InitConfig()

	mux := http.NewServeMux()

	mux.HandleFunc("/put", api.PutFile)
	mux.HandleFunc("/get", api.GetFile)
	mux.HandleFunc("/delete", api.DeleteFile)
	mux.HandleFunc("/upload/status", api.CheckUploadStatus)

	mux.HandleFunc("/metrics", api.GetNetworkMetrics)

	address := fmt.Sprintf(":%s", config.AppConfig.Port)
	log.Printf("[Gateway] Listening on %s", address)

	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("[Main] Server failed: %v", err)
	}
}
