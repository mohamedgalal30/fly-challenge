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

// init registers the provider with the program.
func init() {
	provider := ProviderA{"flaypayA"}
	search.Register(provider.Name, provider)
}

// Search looks at the document for the specified query.
func (p ProviderA) Search(query *search.Query) ([]*search.Result, error) {
	var results []*search.Result
	jsonQuery := gojsonq.New().File(providerAData).From("transactions")

	currency := query.Currency
	statusCode := query.StatusCode
	amountMin := query.AmountMin
	amountMax := query.AmountMax

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
	transactions := jsonQuery.Get()
	transactionsSlice, _ := transactions.([]interface{})

	for _, transaction := range transactionsSlice {
		t := transaction.(map[string]interface{})
		results = append(results, &search.Result{
			Provider:      "flypayA",
			Amount:        t["amount"].(float64),
			Currency:      t["currency"].(string),
			StatusCode:    aStatusString[t["statusCode"].(float64)],
			OrderID:       t["orderReference"].(string),
			TransactionID: t["transactionId"].(string),
		})
	}
	return results, nil
}

// GetName of the provider
func (p ProviderA) GetName() string {
	return p.Name
}
