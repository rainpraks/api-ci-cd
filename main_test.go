package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type MockServer struct {
	server *httptest.Server
}

func SetupMockServer() *MockServer {
	ms := &MockServer{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {

		case r.Method == http.MethodGet && r.URL.Path == "/":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success":true,"data":[{"id":1,"title":"Deal 1"},{"id":2,"title":"Deal 2"}]}`))

		case r.Method == http.MethodGet && r.URL.Path == "/123":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success":true,"data":{"id":123,"title":"Deal 123"}}`))

		case r.Method == http.MethodPost && r.URL.Path == "/":
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"success":true,"data":{"id":999,"title":"New Deal"}}`))

		case r.Method == http.MethodPut && r.URL.Path == "/123":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success":true,"data":{"id":123,"title":"Updated Deal"}}`))

		case r.Method == http.MethodDelete && r.URL.Path == "/123":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success":true,"data":{"id":123,"deleted":true}}`))

		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"success":false,"error":"Not found"}`))
		}
	})

	ms.server = httptest.NewServer(handler)
	return ms
}

func (ms *MockServer) Close() {
	ms.server.Close()
}

func (ms *MockServer) GetURL() string {
	return ms.server.URL
}

func setupTestEnv(t *testing.T) (*mux.Router, *MockServer) {

	mockServer := SetupMockServer()

	os.Setenv("PIPEDRIVE_API_TOKEN", "test-token")
	os.Setenv("PIPEDRIVE_API_URL", mockServer.GetURL())

	ApiToken = os.Getenv("PIPEDRIVE_API_TOKEN")
	PipedriveAPI = os.Getenv("PIPEDRIVE_API_URL")

	router := mux.NewRouter()
	router.HandleFunc("/deals", getDeals).Methods("GET")
	router.HandleFunc("/deals/{id}", getDeals).Methods("GET")
	router.HandleFunc("/deals", postDeal).Methods("POST")
	router.HandleFunc("/deals/{id}", putDeal).Methods("PUT")
	router.HandleFunc("/deals/{id}", deleteDeal).Methods("DELETE")

	return router, mockServer
}

func TestGetDeals(t *testing.T) {
	router, mockServer := setupTestEnv(t)
	defer mockServer.Close()

	req, err := http.NewRequest("GET", "/deals", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.Contains(t, rr.Body.String(), `"id":1,"title":"Deal 1"`)
	assert.Contains(t, rr.Body.String(), `"id":2,"title":"Deal 2"`)
}

func TestGetDealById(t *testing.T) {
	router, mockServer := setupTestEnv(t)
	defer mockServer.Close()

	req, err := http.NewRequest("GET", "/deals/123", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.Contains(t, rr.Body.String(), `"id":123,"title":"Deal 123"`)
}

func TestPostDeal(t *testing.T) {
	router, mockServer := setupTestEnv(t)
	defer mockServer.Close()

	payload := []byte(`{"title":"New Deal","value":1000}`)

	req, err := http.NewRequest("POST", "/deals", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	assert.Contains(t, rr.Body.String(), `"id":999,"title":"New Deal"`)
}

func TestPutDeal(t *testing.T) {
	router, mockServer := setupTestEnv(t)
	defer mockServer.Close()

	payload := []byte(`{"title":"Updated Deal","value":2000}`)

	req, err := http.NewRequest("PUT", "/deals/123", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.Contains(t, rr.Body.String(), `"id":123,"title":"Updated Deal"`)
}

func TestDeleteDeal(t *testing.T) {
	router, mockServer := setupTestEnv(t)
	defer mockServer.Close()

	req, err := http.NewRequest("DELETE", "/deals/123", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.Contains(t, rr.Body.String(), `"id":123,"deleted":true`)
}

func TestErrorHandling(t *testing.T) {
	router, mockServer := setupTestEnv(t)
	defer mockServer.Close()

	errorMockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		conn, _, _ := w.(http.Hijacker).Hijack()
		conn.Close()
	}))
	defer errorMockServer.Close()

	originalURL := PipedriveAPI
	PipedriveAPI = errorMockServer.URL
	defer func() { PipedriveAPI = originalURL }()

	req, _ := http.NewRequest("GET", "/deals", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	req, _ = http.NewRequest("POST", "/deals", bytes.NewBufferString(`{}`))
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	req, _ = http.NewRequest("PUT", "/deals/123", bytes.NewBufferString(`{}`))
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	req, _ = http.NewRequest("DELETE", "/deals/123", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
