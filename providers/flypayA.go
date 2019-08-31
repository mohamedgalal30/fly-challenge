package providers

import (
	"os"

	// "reflect"
	"fly/search"

	"github.com/thedevsaddam/gojsonq"
)

var gopath = os.Getenv("GOPATH")
var providerAData = gopath + "/src/fly/data/flypayA.json"

// ProviderA supply test logic for flypayA provider
type ProviderA struct {
	Name string
}

var aStatus = map[string]int{
	"authorised": 1,
	"decline":    2,
	"refunded":   3,
}
var aStatusString = map[float64]string{
	1: "authorised",
	2: "decline",
	3: "refunded",
}

type aScheme struct {
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
	StatusCode     float64 `json:"statusCode"`
	OrderReference string  `json:"orderReference"`
	TransactionID  string  `json:"transactionId"`
}

// init registers the provider with the program.
func init() {
	provider := ProviderA{"flaypayA"}
	search.Register(provider.Name, provider)
}

// Search looks at the document for the specified query.
func (p ProviderA) Search(query *search.Query) ([]*search.Result, error) {
	var results []*search.Result

	currency := query.Currency
	statusCode := query.StatusCode
	amountMin := query.AmountMin
	amountMax := query.AmountMax

	jsonQuery := gojsonq.New().File(providerAData).From("transactions")
	if currency != "" {
		jsonQuery.WhereEqual("currency", currency)
	}
	if statusCode != "" {
		jsonQuery.WhereEqual("statusCode", aStatus[statusCode])
	}
	if amountMin != 0 || amountMax != 0 {
		jsonQuery.Where("amount", "gte", amountMin)
	}
	if amountMax != 0 {
		jsonQuery.Where("amount", "lte", amountMax)
	}
	var transactionsSlice []aScheme
	jsonQuery.Out(&transactionsSlice)

	for _, transaction := range transactionsSlice {
		results = append(results, &search.Result{
			Provider:      "flypayA",
			Amount:        transaction.Amount,
			Currency:      transaction.Currency,
			StatusCode:    aStatusString[transaction.StatusCode],
			OrderID:       transaction.OrderReference,
			TransactionID: transaction.TransactionID,
		})
	}
	return results, nil
}

// GetName of the provider
func (p ProviderA) GetName() string {
	return p.Name
}
