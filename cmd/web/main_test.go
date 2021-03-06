package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

/*
Used stuff from this site to help me with making tests: https://elithrar.github.io/article/testing-http-handlers-go/
*/

func TestGetLatest(t *testing.T) {
	var data Latest
	data.BaseCurrency = "EUR"
	data.TargetCurrency = "USD"
	m, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/", ioutil.NopCloser(strings.NewReader(string(m))))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getLatest)

	handler.ServeHTTP(rr, req)
}

func TestGetWebhooks(t *testing.T) {
	req, err := http.NewRequest("GET", "/59f9c8a16b022400044b8c74", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getWebhooks)

	handler.ServeHTTP(rr, req)
}

func TestDeleteWebhooks(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/59f9c8a16b022400044b8c74", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteWebhooks)

	handler.ServeHTTP(rr, req)
}

func TestGetAverage(t *testing.T) {
	//Send message
	var data Latest
	data.BaseCurrency = "EUR"
	data.TargetCurrency = "USD"
	m, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/average", ioutil.NopCloser(strings.NewReader(string(m))))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getAverage)

	handler.ServeHTTP(rr, req)
}

func TestRoot2(t *testing.T) { //This one could be giving 5.7% more test coverage if the webhook object the user tries to put in
	//Send message                          //is not in the database from before
	var data Postload
	data.WebhookURL = "https://discordapp.com/api/webhooks/373975976834498560/S9vVxSvLRHpA3V8-F-EAKoB2IGlf0kpUvrJSeYtFI7dzCcCNnkebfiLd0yngTc2UtwF-"
	data.BaseCurrency = "EUR"
	data.TargetCurrency = "USD"
	data.MinTriggerValue = 1
	data.MaxTriggerValue = 3
	m, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/", ioutil.NopCloser(strings.NewReader(string(m))))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(root)

	handler.ServeHTTP(rr, req)
}

func TestRoot(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(root)

	handler.ServeHTTP(rr, req)

}

func TestEvaluationTrigger(t *testing.T) { //Webhook messages get sent everytime ctrl+save
	req, err := http.NewRequest("GET", "/evaluationtrigger", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(triggerwebhooks)

	handler.ServeHTTP(rr, req)

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
