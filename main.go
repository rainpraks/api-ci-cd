package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var apiToken string
var pipedriveAPI string

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system environment variables.")
	}

	apiToken = os.Getenv("PIPEDRIVE_API_TOKEN")
	if apiToken == "" {
		log.Fatal("Missing PIPEDRIVE_API_TOKEN environment variable")
	}
	pipedriveAPI = os.Getenv("PIPEDRIVE_API_URL")
	if pipedriveAPI == "" {
		log.Fatal("Missing PIPEDRIVE_API_URL")
	}
}

func getDeals(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dealID, exists := vars["id"] // Check if an ID is provided

	var url string
	if exists {
		// If an ID is provided, fetch a specific deal
		url = fmt.Sprintf("%s/%s?api_token=%s", pipedriveAPI, dealID, apiToken)
	} else {
		// If no ID is provided, fetch all deals
		url = fmt.Sprintf("%s?api_token=%s", pipedriveAPI, apiToken)
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
	url := fmt.Sprintf("%s?api_token=%s", pipedriveAPI, apiToken)
	resp, err := http.Post(url, "application/json", r.Body)
	if err != nil {
		http.Error(w, "Failed to create deal", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(resp.Body)
}

func putDeal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dealID := vars["id"]

	url := fmt.Sprintf("%s/%s?api_token=%s", pipedriveAPI, dealID, apiToken)

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
	vars := mux.Vars(r)
	dealID := vars["id"]

	url := fmt.Sprintf("%s/%s?api_token=%s", pipedriveAPI, dealID, apiToken)

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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/deals", getDeals).Methods("GET")
	r.HandleFunc("/deals/{id}", getDeals).Methods("GET")
	r.HandleFunc("/deals", postDeal).Methods("POST")
	r.HandleFunc("/deals/{id}", putDeal).Methods("PUT")
	r.HandleFunc("/deals/{id}", deleteDeal).Methods("DELETE")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
