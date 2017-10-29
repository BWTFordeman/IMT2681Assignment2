package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
_id, webhookURL, baseCurrency, targetCurrency, minTriggerValue, maxTriggerValue, currentRate*/

//Person testing webhook stuff.
type Person struct {
	Name string
	Age  int
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
		} else { //Make the request a part of the database
			fmt.Fprintln(w, "Data has been created in database:", p.WebhookURL, p.BaseCurrency, p.TargetCurrency, p.MinTriggerValue, p.MaxTriggerValue)
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

	//TESTING webhook
	request := "https://discordapp.com/api/webhooks/364353373165846528/2Vh8fgXrnsxYQ_MfZuNzW2zwFyN5drj2-wyDo_mUHGIqOiSNWDA-CRx6UmwtqR7D6BhJ"

	form := url.Values{
		"content":    {"xiaoming"},
		"avatar_url": {"https://encrypted-tbn0.gastic.com/images?q=tbn:ANd9GcQ193Rr5CoppqeJDORN3wrGQMIpyWawFAAi2NfjcWccoGd10zL"},
	}

	body := bytes.NewBufferString(form.Encode())
	rsp, err := http.Post(request, "application/json", body)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}
	defer rsp.Body.Close()
	bodyByte, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}
	fmt.Fprintln(w, string(bodyByte))

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
