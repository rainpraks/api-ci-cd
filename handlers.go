package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var ApiToken string

var PipedriveAPI string

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system environment variables.")
	}

	ApiToken = os.Getenv("PIPEDRIVE_API_TOKEN")
	if ApiToken == "" {
		log.Fatal("Missing PIPEDRIVE_API_TOKEN environment variable")
	}
	PipedriveAPI = os.Getenv("PIPEDRIVE_API_URL")
	if ApiToken == "" {
		log.Fatal("Missing PIPEDRIVE_API_TOKEN environment variable")
	}
}

func getDeals(w http.ResponseWriter, r *http.Request) {
	dealID := mux.Vars(r)["id"]

	url := fmt.Sprintf("%s?api_token=%s", PipedriveAPI, ApiToken)
	if dealID != "" {
		url = fmt.Sprintf("%s/%s?api_token=%s", PipedriveAPI, dealID, ApiToken)
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching deals: %v", err)
		http.Error(w, "Failed to fetch deals", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func postDeal(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s?api_token=%s", PipedriveAPI, ApiToken)
	resp, err := http.Post(url, "application/json", r.Body)
	if err != nil {
		http.Error(w, "Failed to create deal", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func putDeal(w http.ResponseWriter, r *http.Request) {
	dealID := mux.Vars(r)["id"]

	url := fmt.Sprintf("%s/%s?api_token=%s", PipedriveAPI, dealID, ApiToken)

	req, err := http.NewRequest(http.MethodPut, url, r.Body)
	if err != nil {
		http.Error(w, "Failed to update deal", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to update deal", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func deleteDeal(w http.ResponseWriter, r *http.Request) {
	dealID := mux.Vars(r)["id"]

	url := fmt.Sprintf("%s/%s?api_token=%s", PipedriveAPI, dealID, ApiToken)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		http.Error(w, "Failed to create delete request", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to delete deal", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
