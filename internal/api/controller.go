package api

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/FraMan97/kairos-p2p-gateway/internal/config"
)

func GetNetworkMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	targetURL := fmt.Sprintf("%s/network/overview", config.AppConfig.KairosExplorerURL)
	resp, err := http.Get(targetURL)
	if err != nil {
		log.Printf("[Metrics] Error contacting explorer: %v", err)
		http.Error(w, "Explorer unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func PutFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("[PutFile] - Only POST method allowed!")
		http.Error(w, "Only POST method allowed!", http.StatusMethodNotAllowed)
		return
	}

	log.Println("[PutFile] - Processing upload request...")

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("[PutFile] - Error retrieving file from form: %v", err)
		http.Error(w, "Error retrieving file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	releaseTime := r.FormValue("release_time")

	bodyReader, bodyWriter := io.Pipe()
	multiWriter := multipart.NewWriter(bodyWriter)

	go func() {
		defer bodyWriter.Close()
		defer multiWriter.Close()

		if err := multiWriter.WriteField("release_time", releaseTime); err != nil {
			log.Printf("[PutFile] - Error writing release_time to pipe: %v", err)
			return
		}

		part, err := multiWriter.CreateFormFile("file", header.Filename)
		if err != nil {
			log.Printf("[PutFile] - Error creating form file in pipe: %v", err)
			return
		}

		if _, err := io.Copy(part, file); err != nil {
			log.Printf("[PutFile] - Error streaming file to pipe: %v", err)
			return
		}
	}()

	targetURL := config.AppConfig.KairosNetworkURL + "/put"
	log.Printf("[PutFile] - Forwarding to Kairos Network: %s", targetURL)

	req, err := http.NewRequest(http.MethodPost, targetURL, bodyReader)
	if err != nil {
		log.Printf("[PutFile] - Error creating proxy request: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", multiWriter.FormDataContentType())

	client := &http.Client{
		Timeout: 0,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[PutFile] - Error contacting Kairos Network: %v", err)
		http.Error(w, "Error contacting storage network", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("[PutFile] - Error copying response to client: %v", err)
	}
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("[GetFile] - Only GET method allowed!")
		http.Error(w, "Only GET method allowed!", http.StatusMethodNotAllowed)
		return
	}

	fileId := r.URL.Query().Get("id")
	if fileId == "" {
		http.Error(w, "[GetFile] - Missing id parameter", http.StatusBadRequest)
		return
	}

	log.Printf("[GetFile] - Proxying request for fileId: %s to Kairos Network...\n", fileId)

	targetURL := fmt.Sprintf("%s/get?id=%s", config.AppConfig.KairosNetworkURL, fileId)
	resp, err := http.Get(targetURL)
	if err != nil {
		log.Printf("[GetFile] - Error connecting to Kairos Network: %v\n", err)
		http.Error(w, "Error connecting to internal storage network", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[GetFile] - Upstream error status: %d\n", resp.StatusCode)
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
		return
	}

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(http.StatusOK)

	bytesCopied, err := io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("[GetFile] - Error streaming file to client: %v\n", err)
		return
	}

	log.Printf("[GetFile] - Successfully streamed %d bytes for fileId: %s\n", bytesCopied, fileId)
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Only DELETE allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("[DeleteFile] - Deleting file request...")

	fileId := r.URL.Query().Get("id")
	if fileId == "" {
		http.Error(w, "[DeleteFile] - Missing id parameter", http.StatusBadRequest)
		return
	}

	targetURL := config.AppConfig.KairosNetworkURL + "/delete?id=" + fileId
	log.Printf("[DeleteFile] - Forwarding to Kairos Network: %s", targetURL)

	req, err := http.NewRequest(http.MethodDelete, targetURL, nil)
	if err != nil {
		http.Error(w, "[DeleteFile] - Internal Error", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "[DeleteFile] - Gateway Error: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func CheckUploadStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed!", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	targetURL := fmt.Sprintf("%s/upload/status?id=%s", config.AppConfig.KairosNetworkURL, id)

	resp, err := http.Get(targetURL)
	if err != nil {
		log.Printf("[CheckStatus] - Error contacting internal network: %v", err)
		http.Error(w, "Error contacting storage network", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("[CheckStatus] - Error streaming response: %v", err)
	}
}
