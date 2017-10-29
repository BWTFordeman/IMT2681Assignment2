package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
using the data:			api.fixer.io/latest?base=EUR;symbols=NOK
{"base": "EUR", "date":"2017-10-27", "rates":{"NOK":9.5348}}*/

//invokeWebhook sends messages through webhooks created in the system
//Must take away lang when database is added, and search through database for names instead.
func invokeWebhook(w http.ResponseWriter, lang [32]string) {
	//May add validation for /slack or /github at end of webhookURL
	//Discord has content, slack has text
	webhookURL := "https://discordapp.com/api/webhooks/373975976834498560/S9vVxSvLRHpA3V8-F-EAKoB2IGlf0kpUvrJSeYtFI7dzCcCNnkebfiLd0yngTc2UtwF-"

	res, err := http.PostForm(webhookURL, url.Values{"content": {"baseCurrency: " + lang[0]}, "username": {"CurrencyChecker"}})
	if err != nil {
		fmt.Println(err.Error(), "Panic or something")
	}

	if res.StatusCode == 200 || res.StatusCode == 204 {
		fmt.Fprintln(w, "statuscode: ", res.StatusCode)
	} else {
		fmt.Fprintln(w, "Wrong status: ", res.StatusCode)
	}
}

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

	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var p Postload
		err := decoder.Decode(&p)

		//Check if currencies are of valid types.
		var base = false
		var target = false
		for i := 0; i < len(lang); i++ {
			if p.BaseCurrency == lang[i] {
				base = true
			}
			if p.TargetCurrency == lang[i] {
				target = true
			}
		}

		//Create row in databse if valid:
		if err != nil || base != true || target != true {
			http.Error(w, "Invalid post value", http.StatusBadRequest)
		} else {
			//put in database here.
			fmt.Fprintln(w, "Id for your webhook: ") // add an id generated here.
		}

		defer r.Body.Close()
	} else { //If not post:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	//Webhook:		//Needs to be in another handler
	//invokeWebhook(w, lang)

}

/*//GET AND DELETE WILL BE IN ANOTHER HANDLER
switch r.Method {

case "GET": //Show stuff from database
	fmt.Fprintln(w, "GET")

case "DELETE": //Delete something from the database
}*/
