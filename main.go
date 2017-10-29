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

/*
In the databse these data will be stored:
__id, webhookURL, baseCurrency, targetCurrency, minTriggerValue, maxTriggerValue, currentRate*/

//Message send stuff through webhook to discord
type Message struct {
	Content  string `json:"content"`
	Username string `json:"username"`
}

//Content sent through the Message
/*type Content struct {
	BaseCurrency    string  `json:"baseCurrency"`
	TargetCurrency  string  `json:"targetCurrency"`
	CurrentRate     float32 `json:"currentRate"`
	MinTriggerValue float32 `json:"minTriggerValue"`
	MaxTriggerValue float32 `json:"maxTriggerValue"`
}*/

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
		var target = false
		for i := 0; i < len(lang); i++ {
			if p.BaseCurrency == lang[i] {
				base = true
			}
			if p.TargetCurrency == lang[i] {
				target = true
			}
		}

		if err != nil || base != true || target != true {
			http.Error(w, "Invalid post value", http.StatusBadRequest)
		} else {
			fmt.Fprintln(w, "Data has been created in database:", p.WebhookURL, p.BaseCurrency, p.TargetCurrency, p.MinTriggerValue, p.MaxTriggerValue) //Do not use this!
			//put in database here.
			fmt.Fprintln(w, "webhook id generated during registration: ") // add an id generated here.
		}
		defer r.Body.Close()
		//GET AND DELETE WILL BE IN ANOTHER HANDLER
	case "GET": //Show stuff from database
		fmt.Fprintln(w, "GET")

	case "DELETE": //Delete something from the database

	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
	//Add a function that runs every 24 hour.

	//TESTING webhook						/May add validation for /slack or /github at end of webhookURL
	//Discord has content, slack has text
	webhookURL := "https://discordapp.com/api/webhooks/373975976834498560/S9vVxSvLRHpA3V8-F-EAKoB2IGlf0kpUvrJSeYtFI7dzCcCNnkebfiLd0yngTc2UtwF-"

	var message Message
	message.Username = "Fordeman"
	message.Content = (`{\n"baseCurrency": ` + `NOK\n}`)

	/*msh, err := json.Marshal(message)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}*/

	//res, err := http.Post(webhookURL, "application/json", bytes.NewReader(msh))
	res, err := http.PostForm(webhookURL, url.Values{"content": {"bar"}, "username": {"testName"}})
	if err != nil {
		fmt.Println(err.Error(), "Panic or something")
	}

	if res.StatusCode == 200 || res.StatusCode == 204 {
		fmt.Fprintln(w, "statuscode: ", res.StatusCode)
	} else {
		fmt.Fprintln(w, "Wrong status: ", res.StatusCode)
	}

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
