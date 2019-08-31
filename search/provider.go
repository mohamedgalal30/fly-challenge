package search

import (
	"log"
)

// Result contains the result of a search.
type Result struct {
	Provider      string
	Amount        float64
	Currency      string
	StatusCode    string
	OrderID       string
	TransactionID string
}

// Provider defines the behavior required by types that want
// to implement a new provider type.
type Provider interface {
	GetName() string
	Search(query *Query) ([]*Result, error)
}

// Match is launched as a goroutine for each individual provider to run
// searches concurrently.
func Match(provider Provider, query *Query, results chan<- *Result) {
	// Perform the search against the specified provider.
	searchResults, err := provider.Search(query)
	if err != nil {
		log.Println(err)
		return
	}

	// Write the results to the channel.
	for _, result := range searchResults {
		results <- result
	}
}
