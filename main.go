package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.Use(requestLogger)
	r.HandleFunc("/deals", getDeals).Methods("GET")
	r.HandleFunc("/deals/{id}", getDeals).Methods("GET")
	r.HandleFunc("/deals", postDeal).Methods("POST")
	r.HandleFunc("/deals/{id}", putDeal).Methods("PUT")
	r.HandleFunc("/deals/{id}", deleteDeal).Methods("DELETE")
	r.HandleFunc("/metrics", getMetrics).Methods("GET")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
