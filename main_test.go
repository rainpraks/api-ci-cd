package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Mock environment variable
	os.Setenv("PIPEDRIVE_API_TOKEN", "test_token")
	exitVal := m.Run()
	os.Exit(exitVal)
}

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/deals", getDeals).Methods("GET")
	r.HandleFunc("/deals/{id}", getDeals).Methods("GET")
	r.HandleFunc("/deals", postDeal).Methods("POST")
	r.HandleFunc("/deals/{id}", putDeal).Methods("PUT")
	r.HandleFunc("/deals/{id}", deleteDeal).Methods("DELETE")
	r.HandleFunc("/metrics", getMetrics).Methods("GET")
	return r
}

func TestGetDeals(t *testing.T) {
	req, err := http.NewRequest("GET", "/deals", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), `"success":true`)
}

func TestGetDealByID(t *testing.T) {
	req, err := http.NewRequest("GET", "/deals/123", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestPostDeal(t *testing.T) {
	deal := map[string]string{"title": "New Deal"}
	dealJSON, _ := json.Marshal(deal)

	req, err := http.NewRequest("POST", "/deals", bytes.NewBuffer(dealJSON))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code) // Expecting 201 Created
}

func TestPutDeal(t *testing.T) {
	deal := map[string]string{"title": "Updated Deal"}
	dealJSON, _ := json.Marshal(deal)

	req, err := http.NewRequest("PUT", "/deals/123", bytes.NewBuffer(dealJSON))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	assert.Contains(t, []int{http.StatusOK, http.StatusNotFound}, rr.Code) // Either 200 or 404
}

func TestDeleteDeal(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/deals/123", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	assert.Contains(t, []int{http.StatusOK, http.StatusGone}, rr.Code) // Either 200 or 410
}

func TestGetMetrics(t *testing.T) {
	req, err := http.NewRequest("GET", "/metrics", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), `"total_requests"`)
}
