package providers

import (
	"fly/search"

	"github.com/thedevsaddam/gojsonq"
)

var providerBData = gopath + "/src/fly/data/flypayB.json"

// ProviderB supply search logic for flypayB provider
type ProviderB struct {
	Name string
}

type bScheme struct {
	Value               float64 `json:"value"`
	TransactionCurrency string  `json:"transactionCurrency"`
	StatusCode          float64 `json:"statusCode"`
	OrderInfo           string  `json:"orderInfo"`
	PaymentID           string  `json:"paymentID"`
}

var bStatus = map[string]int{
	"authorised": 100,
	"decline":    200,
	"refunded":   300,
}

var bStatusString = map[float64]string{
	100: "authorised",
	200: "decline",
	300: "refunded",
}

// init registers the provider with the program.
func init() {
	var provider ProviderB
	search.Register("flypayB", provider)
}

// Search looks at the document for the specified query.
func (p ProviderB) Search(query *search.Query) ([]*search.Result, error) {
	var results []*search.Result

	jsonQuery := gojsonq.New().File(providerBData).From("transactions")

	currency := query.Currency
	statusCode := query.StatusCode
	amountMin := query.AmountMin
	amountMax := query.AmountMax

	if currency != "" {
		jsonQuery.WhereEqual("transactionCurrency", currency)
	}
	if statusCode != "" {
		jsonQuery.WhereEqual("statusCode", bStatus[statusCode])
	}
	if amountMin != 0 || amountMax != 0 {
		jsonQuery.Where("value", "gte", amountMin)
	}
	if amountMax != 0 {
		jsonQuery.Where("value", "lte", amountMax)
	}
	var transactionsSlice []bScheme
	jsonQuery.Out(&transactionsSlice)
	for _, transaction := range transactionsSlice {
		results = append(results, &search.Result{
			Provider:      "flypayB",
			Amount:        transaction.Value,
			Currency:      transaction.TransactionCurrency,
			StatusCode:    bStatusString[transaction.StatusCode],
			OrderID:       transaction.OrderInfo,
			TransactionID: transaction.PaymentID,
		})
	}
	return results, nil
}

// GetName of the provider
func (p ProviderB) GetName() string {
	return p.Name
}
