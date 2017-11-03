package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoot(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(root)

	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

/*func TestRoot2(t *testing.T) {
	req, err := http.NewRequest("POST", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(root)

	handler.ServeHTTP(rr, req)
}*/

func TestFindAllWebhooks(t *testing.T) {
	webhooks, err := findAllWebhooks()
	if err != nil {
		t.Errorf("Could not find any webhooks %v", webhooks)
	}
}

func TestGetCurrentValue(t *testing.T) {
	f := Fixer{}
	f.Rates = map[string]float64{"AUD": 1.5117, "BGN": 1.9558, "BRL": 3.8047, "CAD": 1.496, "CHF": 1.1647, "CNY": 7.7011, "CZK": 25.535, "DKK": 7.4418, "GBP": 0.8869, "HKD": 9.0851, "HRK": 7.5302,
		"HUF": 310.9, "IDR": 15753.0, "ILS": 4.0845, "INR": 75.229, "JPY": 132.9, "KRW": 1294.9, "MXN": 22.232, "MYR": 4.9264, "NOK": 9.4838, "NZD": 1.6867, "PHP": 59.936, "PLN": 4.2376,
		"RON": 4.5992, "RUB": 68.033, "SEK": 9.7615, "SGD": 1.5843, "THB": 38.568, "TRY": 4.4519, "USD": 1.1645, "ZAR": 16.288}
	i := getCurrentValue(f, "ZAR")
	if i == 0 {
		t.Errorf("Expected a value for ZAR, got:%v", i)
	}
	j := getCurrentValue(f, "NOO")
	if i == 0 {
		t.Errorf("Expected a value for NOO, got:%v", j)
	}
}
