package search_test

import (
	_ "fly/providers"
	"fly/search"
	"reflect"
	"strconv"
	"testing"
)

const checkMark = "\033[32m" + "\u2713" + "\033[39m"
const ballotX = "\033[31m" + "\u2717" + "\033[39m"

func TestRun(t *testing.T) {
	t.Log("Given the need to retrieve the providers transactions that match the query.")
	{
		testFindAllTransactions(t)
		testFilterTransactionsByProvider(t)
		testFilterTransactionsByStatus(t)
		testCombinedFilterTransactions(t)
	}
}

func testFindAllTransactions(t *testing.T) {
	t.Log("\tWhen try to find all transactions")
	{
		allQuery := search.Query{}
		result := search.Run(&allQuery)
		if len(result) == 15 {
			t.Log("\t\tShould retrieve all available data.", checkMark)
		} else {
			t.Error("\t\tShould retrieve all available data.", ballotX)
		}
	}

}

func testFilterTransactionsByProvider(t *testing.T) {
	t.Logf("\tWhen try to filter transactions by provider")
	{
		providerQueries := []*search.Query{
			{Provider: "flypayA"},
			{Provider: "flypayB"},
		}
		for _, providerQuery := range providerQueries {
			providerName := providerQuery.Provider
			result := search.Run(providerQuery)
			providerQueryPassed := true
			for _, transaction := range result {
				if transaction.Provider != providerName {
					providerQueryPassed = false
					break
				}
			}
			if providerQueryPassed {
				t.Log("\t\tShould retrieve provider \""+providerName+"\" data.", checkMark)
			} else {
				t.Error("\t\tShould retrieve provider \""+providerName+"\" data.", ballotX)
			}
		}
	}
}
func testFilterTransactionsByStatus(t *testing.T) {
	t.Logf("\tWhen try to filter transactions by status code")
	{
		statusQueries := []*search.Query{
			{StatusCode: "authorised"},
			{StatusCode: "decline"},
			{StatusCode: "refunded"},
		}
		for _, statusQuery := range statusQueries {
			status := statusQuery.StatusCode
			statusQueryPassed := true
			result := search.Run(statusQuery)
			for _, transaction := range result {
				if transaction.StatusCode != status {
					statusQueryPassed = false
					break
				}
			}
			if statusQueryPassed {
				t.Log("\t\tShould retrieve status  \""+status+"\" data.", checkMark)
			} else {
				t.Error("\t\tShould retrieve status  \""+status+"\" data.", ballotX)
			}
		}
	}
}
func testCombinedFilterTransactions(t *testing.T) {
	t.Logf("\tWhen try to filter transactions by combined query")
	{
		combinedQueries := []search.Query{
			{Provider: "flypayA", AmountMin: 500},
			{Provider: "flypayB", AmountMax: 1000, StatusCode: "authorised"},
			{Provider: "flypayA", Currency: "AUD"},
		}
		for queryIndex, combinedQuery := range combinedQueries {
			queryType := reflect.TypeOf(combinedQuery)
			queryValues := reflect.ValueOf(combinedQuery)
			combindeQueryPassed := true
			result := search.Run(&combinedQuery)

			for _, transaction := range result {
				transactionValues := reflect.ValueOf(transaction)
				for i := 0; i < queryValues.NumField(); i++ {
					queryFieldName := queryType.Field(i).Name
					queryFieldValue := queryValues.Field(i)

					transactionFieldValue := reflect.Indirect(transactionValues).FieldByName(queryFieldName)
					if queryFieldName == "AmountMin" || queryFieldName == "AmountMax" {
						if queryFieldValue.Interface().(float64) == 0 {
							continue
						}
						transactionFieldValue = reflect.Indirect(transactionValues).FieldByName("Amount")
						passed := (queryFieldName == "AmountMin" &&
							transactionFieldValue.Interface().(float64) >= queryFieldValue.Interface().(float64)) ||
							(queryFieldName == "AmountMax" &&
								transactionFieldValue.Interface().(float64) <= queryFieldValue.Interface().(float64))
						if passed {

							continue
						} else {
							combindeQueryPassed = false
							break
						}
					}

					if transactionFieldValue.Interface() != queryFieldValue.Interface() {
						if queryFieldValue.Interface().(string) == "" {
							continue
						}

						combindeQueryPassed = false
						break
					}
				}
			}
			numstr := strconv.Itoa(queryIndex + 1)
			if combindeQueryPassed {
				t.Log("\t\tShould retrieve the data which match the combined query number "+numstr+".", checkMark)
			} else {
				t.Error("\t\tShould retrieve the data which match the combined query number "+numstr+".", ballotX)
			}
		}
	}
}

func TestRegister(t *testing.T) {

	t.Log("Given the need to test Register Provider.")
	{
		t.Logf("\tWhen try to register new provider")
		{
			var provider search.Provider
			providerName := "flyProvider"
			search.Register(providerName, provider)
			_, found := search.Providers[providerName]
			if !found {
				t.Fatal("\t\tShould found the registered provider in providers map.", ballotX)
			}
			t.Log("\t\tShould found the registered provider in providers map.", checkMark)
		}
	}
}
