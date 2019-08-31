package search_test

import (
	_ "fly/providers"
	"fly/search"
	"reflect"
	"testing"
)

const checkMark = "\033[32m" + "\u2713" + "\033[39m"
const ballotX = "\033[31m" + "\u2717" + "\033[39m"

func TestRun(t *testing.T) {
	t.Log("Given the need to retrieve the providers transactions that match the query.")
	{
		t.Log("\tWhen try to find all transactions")
		{
			allQuery := search.Query{}
			result := search.Run(allQuery)
			if len(result) == 15 {
				t.Log("\t\tShould retrieve all available data.", checkMark)
			} else {
				t.Error("\t\tShould retrieve all available data.", ballotX)
			}
		}
		t.Logf("\tWhen try to filter transactions by provider")
		{
			providerQueries := []*search.Query{
				{Provider: "flypayA"},
				{Provider: "flypayB"},
			}
			for _, providerQuery := range providerQueries {
				providerName := providerQuery.Provider
				result := search.Run(*providerQuery)
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
				result := search.Run(*statusQuery)
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
		t.Logf("\tWhen try to filter transactions by combined query")
		{
			combinedQueries := []search.Query{
				{
					Provider:   "flypayA",
					StatusCode: "authorised",
					AmountMin:  500,
					AmountMax:  1500,
					Currency:   "",
				},
			}
			for _, combinedQuery := range combinedQueries {
				fields := reflect.TypeOf(combinedQuery)
				values := reflect.ValueOf(combinedQuery)

				// values := make([]interface{}, v.NumField())
				// t.Log(values)
				// for i := 0; i < v.NumField(); i++ {
				// 	t.Log(v.MapKey(i).Interface())
				// 	field := fields.Field(i)
				// 	value := values.Field(i)
				// 	values[i] = v.Field(i).Interface()
				// }

				t.Log(fields.Field(1).Name)
				t.Log(values)
				// for key, value := range combinedQuery {

				// }
				status := combinedQuery.StatusCode
				combindeQueryPassed := true
				result := search.Run(combinedQuery)
				for _, transaction := range result {
					if transaction.StatusCode != status {
						combindeQueryPassed = false
						break
					}
				}
				if combindeQueryPassed {
					t.Log("\t\tShould retrieve status  \""+status+"\" data.", checkMark)
				} else {
					t.Error("\t\tShould retrieve status  \""+status+"\" data.", ballotX)
				}
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
