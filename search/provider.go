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
// to implement a new search type.
type Provider interface {
	Search(query Query) ([]*Result, error)
}

// Match is launched as a goroutine for each individual transaction to run
// searches concurrently.
func Match(provider Provider, query Query, results chan<- *Result) {
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
