package config

import (
	"log"
	"os"
)

type Config struct {
	KairosNetworkURL  string
	KairosExplorerURL string
	KairosLedgerURL   string
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

	ledgerURL := os.Getenv("KAIROS_LEDGER_URL")
	if ledgerURL == "" {
		ledgerURL = "http://kairos-ledger:80"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	AppConfig = &Config{
		KairosNetworkURL:  kairosURL,
		KairosExplorerURL: explorerURL,
		KairosLedgerURL:   ledgerURL,
		Port:              port,
	}

	log.Printf("[Config] Loaded: Network=%s, Explorer=%s, Ledger=%s, Port=%s",
		AppConfig.KairosNetworkURL, AppConfig.KairosExplorerURL, AppConfig.KairosLedgerURL, AppConfig.Port)
}
