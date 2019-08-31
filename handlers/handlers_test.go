// Sample test to show how to test the execution of an
// internal endpoint.
package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"fly/handlers"
	_ "fly/providers"
	"fly/search"
)

const checkMark = "\033[32m" + "\u2713" + "\033[39m"
const ballotX = "\033[31m" + "\u2717" + "\033[39m"

func init() {
	handlers.Routes()
}

// TestSendJSON testing the sendjson internal endpoint.
func TestSearchHandler(t *testing.T) {
	t.Log("Given the need to test the SearchHandler endpoint.")
	{
		url := "/api/payment/transaction"
		t.Logf("\tWhen checking \"%s\" ", url)
		{
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatal("\t\tShould be able to create a request.", ballotX, err)
			}
			t.Log("\t\tShould be able to create a request.", checkMark)

			w := httptest.NewRecorder()
			http.DefaultServeMux.ServeHTTP(w, req)

			if w.Code != 200 {
				t.Fatal("\t\tShould receive \"200\"", ballotX, w.Code)
			}
			t.Log("\t\tShould receive \"200\"", checkMark)

			result := []search.Result{}
			if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
				t.Fatal("\t\tShould decode the response.", ballotX)
			}
			t.Log("\t\tShould decode the response.", checkMark)

			if len(result) == 15 {
				t.Log("\t\tShould have data array of length of 15 .", checkMark)
			} else {
				t.Error("\t\tShould have data array of length of 15 .", ballotX)
			}
		}
	}
}
