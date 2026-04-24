package config

import (
	"log"
	"os"
)

type Config struct {
	KairosNetworkURL  string
	KairosExplorerURL string
	Port              string
}

var AppConfig *Config

func InitConfig() {
	kairosURL := os.Getenv("KAIROS_NETWORK_URL")
	if kairosURL == "" {
		kairosURL = "http://kairos-engine-api:80"
	}

	explorerURL := os.Getenv("KAIROS_EXPLORER_URL")
	if explorerURL == "" {
		explorerURL = "http://kairos-explorer:8081"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	AppConfig = &Config{
		KairosNetworkURL:  kairosURL,
		KairosExplorerURL: explorerURL,
		Port:              port,
	}

	log.Printf("[Config] Loaded: Network=%s, Explorer=%s, Port=%s",
		AppConfig.KairosNetworkURL, AppConfig.KairosExplorerURL, AppConfig.Port)
}
