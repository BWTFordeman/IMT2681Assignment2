package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

//Postload datatype I get from fixer.io
type Postload struct {
	WebhookURL      string  `json:"webhookURL"`
	BaseCurrency    string  `json:"baseCurrency"`
	TargetCurrency  string  `json:"targetCurrency"`
	MinTriggerValue float32 `json:"minTriggerValue"`
	MaxTriggerValue float32 `json:"maxTriggerValue"`
}

/*
type Payload struct {
	WebhookURL   string `json:"url"`
	BaseCurrency string `json:"baseCurrency"`
	Date         int    `json:"date"`
	Rates        struct {
		currency map[string]int
	} `json:"rates"`
}*/

//Make another one for requesting data from fixer.io
/*
using the data:			fixer.io/latest?base=EUR;symbols=NOK
{"base": "EUR", "date":"2017-10-27", "rates":{"NOK":9.5348}}*/

/*
In the databse these data will be stored:
webhookURL, baseCurrency, targetCurrency, minTriggerValue, maxTriggerValue, currentRate*/

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		fmt.Println(err.Error(), "Panic or something")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	lang := [...]string{"EUR", "AUD", "BGN", "BRL", "CAD", "CHF", "CNY", "CZK", "DKK", "GBP", "HKD", "HRK", "HUF", "IDR", "ILS", "INR", "JPY", "KRW",
		"MXN", "MYR", "NOK", "NZD", "PHP", "PLN", "RON", "RUB", "SEK", "SGD", "THB", "TRY", "USD", "ZAR"}

	//Goto right service:
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		var p Postload
		err := decoder.Decode(&p)

		//Check if currencies are of valid types.
		var base = false
		for i := 0; i < len(lang); i++ {
			if p.BaseCurrency == lang[i] {
				base = true
			}
		}

		if err != nil && base == true {
			http.Error(w, "Invalid post value", http.StatusBadRequest)
		} else {
			fmt.Fprintln(w, "Data has been created in database:", p.WebhookURL, p.BaseCurrency, p.TargetCurrency, p.MinTriggerValue, p.MaxTriggerValue)
			//Create p."data" in the database.
		}
		defer r.Body.Close()

	case "GET":
		fmt.Fprintln(w, "GET")

	case "DELETE":

	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
	//Add a function that runs every 24 hour.

	fmt.Fprintln(w, "Daily func")
	/*
		timer := time.NewTimer(time.Hour * 24)
	*/

	/*go func() {
		time.Sleep(time.Second * 10)
		timerFinished <- true
	}()

	for {
		select {
		case <-timer:
			fmt.Println("Timer expired")
		case <-timerFinished:
			fmt.Println("Done")
			http.HandleFunc("/", handler)
		}
	}*/
}
