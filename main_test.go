package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	os.Setenv("PIPEDRIVE_API_TOKEN", "mock_token")
	ApiToken = os.Getenv("PIPEDRIVE_API_TOKEN")

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGetDeals(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockResponse := `{"success": true, "data": [{"id": 1, "title": "Test Deal"}]}`
	apiURL := fmt.Sprintf("%s?api_token=%s", PipedriveAPI, ApiToken)

	httpmock.RegisterResponder("GET", apiURL,
		httpmock.NewStringResponder(200, mockResponse))

	req := httptest.NewRequest(http.MethodGet, "/deals", nil)
	rr := httptest.NewRecorder()

	getDeals(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, mockResponse, rr.Body.String())
}

func TestGetDealByID(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	dealID := "123"
	mockResponse := `{"success": true, "data": {"id": 123, "title": "Specific Deal"}}`
	apiURL := fmt.Sprintf("%s/%s?api_token=%s", PipedriveAPI, dealID, ApiToken)

	httpmock.RegisterResponder("GET", apiURL,
		httpmock.NewStringResponder(200, mockResponse))

	req := httptest.NewRequest(http.MethodGet, "/deals/123", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/deals/{id}", getDeals)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, mockResponse, rr.Body.String())
}

func TestPostDeal(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockResponse := `{"success": true, "data": {"id": 456, "title": "New Deal"}}`
	apiURL := fmt.Sprintf("%s?api_token=%s", PipedriveAPI, ApiToken)

	httpmock.RegisterResponder("POST", apiURL,
		httpmock.NewStringResponder(201, mockResponse))

	payload := `{"title": "New Deal"}`
	req := httptest.NewRequest(http.MethodPost, "/deals", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	postDeal(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var expected, actual map[string]interface{}
	json.Unmarshal([]byte(mockResponse), &expected)
	json.Unmarshal(rr.Body.Bytes(), &actual)

	assert.Equal(t, expected, actual)
}

func TestPutDeal(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	dealID := "456"
	mockResponse := `{"success": true, "data": {"id": 456, "title": "Updated Deal"}}`
	apiURL := fmt.Sprintf("%s/%s?api_token=%s", PipedriveAPI, dealID, ApiToken)

	httpmock.RegisterResponder("PUT", apiURL,
		httpmock.NewStringResponder(200, mockResponse))

	payload := `{"title": "Updated Deal"}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/deals/%s", dealID), bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/deals/{id}", putDeal)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, mockResponse, rr.Body.String())
}

func TestDeleteDeal(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	dealID := "789"
	mockResponse := `{"success": true, "data": {"id": 789}}`
	apiURL := fmt.Sprintf("%s/%s?api_token=%s", PipedriveAPI, dealID, ApiToken)

	httpmock.RegisterResponder("DELETE", apiURL,
		httpmock.NewStringResponder(200, mockResponse))

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/deals/%s", dealID), nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/deals/{id}", deleteDeal)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, mockResponse, rr.Body.String())
}
