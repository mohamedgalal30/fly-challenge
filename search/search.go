package search

import (
	"log"
	"sync"
)

//Providers is a map of registered providers for searching.
var Providers = make(map[string]Provider)

// Query  type to recieve and permit query params
type Query struct {
	Provider   string
	StatusCode string
	AmountMin  float64
	AmountMax  float64
	Currency   string
}

// Run performs the search logic.
func Run(query *Query) []*Result {

	// Create an unbuffered channel to receive match results to display.
	results := make(chan *Result)
	searchProviders := getSearchProviders(query.Provider)
	// Setup a wait group so we can process all providers.
	var waitGroup sync.WaitGroup

	// Set the number of goroutines we need to wait for while they process the providers.
	waitGroup.Add(len(searchProviders))

	// Launch a goroutine for each provider to find the results.
	for _, provider := range searchProviders {

		// Launch the goroutine to perform the search.
		go func(provider Provider) {
			Match(provider, query, results)
			waitGroup.Done()
		}(provider)
	}

	// Launch a goroutine to monitor when all the work is done.
	go func() {
		// Wait for everything to be processed.
		waitGroup.Wait()

		// Close the channel to signal to the resultChanToSlice function that we can go on.
		close(results)
	}()

	return resultChanToSlice(results)
}

// Register is called to register a provider for use by the program.
func Register(providerName string, provider Provider) {
	if _, exists := Providers[providerName]; exists {
		log.Fatalln(providerName, "Provider already registered")
	}

	log.Println("Register", providerName, "provider")
	Providers[providerName] = provider
}

func getSearchProviders(providerName string) map[string]Provider {
	searchProviders := make(map[string]Provider)
	if providerName != "" {
		provider, found := Providers[providerName]
		if found {
			searchProviders[providerName] = provider
		}
	} else {
		searchProviders = Providers
	}
	return searchProviders
}

func resultChanToSlice(resultCh chan *Result) []*Result {
	var results []*Result
	for result := range resultCh {
		results = append(results, result)
	}
	return results
}
