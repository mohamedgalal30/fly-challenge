// Package handlers provides the endpoints for the web service.
package handlers

import (
	"encoding/json"
	"fly/search"
	"net/http"
	"strconv"
)

// Routes sets the routes for the web service.
func Routes() {
	http.HandleFunc("/api/payment/transaction", SearchHandler)
}

// SearchHandler provide support for search in providers transactions
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	min, _ := strconv.ParseFloat(query.Get("amountMin"), 64)
	max, _ := strconv.ParseFloat(query.Get("amountMax"), 64)
	q := search.Query{
		Provider:   query.Get("provider"),
		StatusCode: query.Get("statusCode"),
		AmountMin:  min,
		AmountMax:  max,
		Currency:   query.Get("currency"),
	}

	result := search.Run(&q)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&result)
}
