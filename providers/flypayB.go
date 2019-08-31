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
	transactions := jsonQuery.Get()
	transactionsSlice, _ := transactions.([]interface{})

	for _, transaction := range transactionsSlice {
		t := transaction.(map[string]interface{})
		results = append(results, &search.Result{
			Provider:      "flypayB",
			Amount:        t["value"].(float64),
			Currency:      t["transactionCurrency"].(string),
			StatusCode:    bStatusString[t["statusCode"].(float64)],
			OrderID:       t["orderInfo"].(string),
			TransactionID: t["paymentId"].(string),
		})
	}
	return results, nil
}

// GetName of the provider
func (p ProviderB) GetName() string {
	return p.Name
}
