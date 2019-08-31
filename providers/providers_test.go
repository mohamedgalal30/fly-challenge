package providers_test

import (
	"fly/providers"
	"fly/search"
	"strconv"
	"testing"
)

const checkMark = "\033[32m" + "\u2713" + "\033[39m"
const ballotX = "\033[31m" + "\u2717" + "\033[39m"

func TestFlypayASearch(t *testing.T) {
	provider := providers.ProviderA{"flypayA"}
	t.Logf("Given the need to test %s search logic.", provider.GetName())
	{
		allTransactionsNumber := 9
		testFindAllTransactions(t, provider, allTransactionsNumber)

		statusTransactionsNum := map[string]int{
			"authorised": 4,
			"decline":    2,
			"refunded":   3,
		}
		testFilterTransactionsByStatus(t, provider, &statusTransactionsNum)

		transactionsNum := []int{6, 4, 3}
		testCombinedFilterTransactions(t, provider, &transactionsNum)
	}
}
func TestFlypayBSearch(t *testing.T) {
	provider := providers.ProviderB{"flypayB"}
	t.Logf("Given the need to test %s search logic.", provider.GetName())
	{
		allTransactionsNumber := 6
		testFindAllTransactions(t, provider, allTransactionsNumber)

		statusTransactionsNum := map[string]int{
			"authorised": 2,
			"decline":    2,
			"refunded":   2,
		}
		testFilterTransactionsByStatus(t, provider, &statusTransactionsNum)

		transactionsNum := []int{2, 1, 2}
		testCombinedFilterTransactions(t, provider, &transactionsNum)
	}
}

func testFindAllTransactions(t *testing.T, provider search.Provider, allTransactionsNumber int) {
	t.Log("\tWhen try to find all transactions of ", provider.GetName())
	{
		allQuery := search.Query{}
		result, _ := provider.Search(&allQuery)
		if len(result) == allTransactionsNumber {
			t.Logf("\t\tShould retrieve all %s data. "+checkMark, provider.GetName())
		} else {
			t.Errorf("\t\tShould retrieve all %s data. "+ballotX, provider.GetName())
		}
	}

}

func testFilterTransactionsByStatus(t *testing.T, provider search.Provider, transactionsNum *map[string]int) {
	t.Logf("\tWhen try to filter transactions by status code")
	{
		statusQueries := []*search.Query{
			{StatusCode: "authorised"},
			{StatusCode: "decline"},
			{StatusCode: "refunded"},
		}

		for _, statusQuery := range statusQueries {
			status := statusQuery.StatusCode
			dataNum := (*transactionsNum)[status]
			result, _ := provider.Search(statusQuery)
			if len(result) == dataNum {
				t.Log("\t\tShould retrieve status  \""+status+"\" data.", checkMark)
			} else {
				t.Error("\t\tShould retrieve status  \""+status+"\" data.", ballotX)
			}
		}
	}
}
func testCombinedFilterTransactions(t *testing.T, provider search.Provider, transactionsNum *[]int) {
	t.Logf("\tWhen try to filter transactions by combined query")
	{
		combinedQueries := []*search.Query{
			{AmountMin: 500, AmountMax: 2000},
			{AmountMax: 1000, StatusCode: "authorised"},
			{Currency: "AUD", StatusCode: "refunded"},
		}

		for queryIndex, combinedQuery := range combinedQueries {
			dataNum := (*transactionsNum)[queryIndex]
			result, _ := provider.Search(combinedQuery)
			numstr := strconv.Itoa(queryIndex + 1)
			if len(result) == dataNum {
				t.Log("\t\tShould retrieve the data which match the combined query number "+numstr+".", checkMark)
			} else {
				t.Error("\t\tShould retrieve the data which match the combined query number "+numstr+".", ballotX)
			}
		}
	}
}
