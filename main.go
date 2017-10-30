package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Postload data retrieved from adding webhook
type Postload struct {
	WebhookURL      string  `json:"webhookURL"`
	BaseCurrency    string  `json:"baseCurrency"`
	TargetCurrency  string  `json:"targetCurrency"`
	MinTriggerValue float32 `json:"minTriggerValue"`
	MaxTriggerValue float32 `json:"maxTriggerValue"`
}

//Webhook retrieves data from the webhook collection:
type Webhook struct {
	WebhookURL      string  `json:"webhookURL"`
	BaseCurrency    string  `json:"baseCurrency"`
	TargetCurrency  string  `json:"targetCurrency"`
	MinTriggerValue float32 `json:"minTriggerValue"`
	MaxTriggerValue float32 `json:"maxTriggerValue"`
	CurrentRate     float32 `json:"currentRate"`
}

//Fixer retrieves latest data from fixer.io							api.fixer.io/latest?base=EUR;symbols=NOK
type Fixer struct {
	BaseCurrency string `json:"baseCurrency"`
	Date         string `json:"date"`
	Rates        struct {
		string float32
	} `json:"rates"`
}

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
		fmt.Fprintln(w, "Wrong status: ", res.StatusCode, http.StatusText(res.StatusCode))
	}
}

func main() {
	http.HandleFunc("/", root)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		fmt.Println(err.Error(), "Panic or something")
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	lang := [...]string{"EUR", "AUD", "BGN", "BRL", "CAD", "CHF", "CNY", "CZK", "DKK", "GBP", "HKD", "HRK", "HUF", "IDR", "ILS", "INR", "JPY", "KRW",
		"MXN", "MYR", "NOK", "NZD", "PHP", "PLN", "RON", "RUB", "SEK", "SGD", "THB", "TRY", "USD", "ZAR"}
	USER := os.Getenv("DB_USER")
	PASSWORD := os.Getenv("DB_PASSWORD")
	tempstring := ("mongodb://" + USER + ":" + PASSWORD + "@ds241055.mlab.com:41055/imt2681")
	session, err := mgo.Dial(tempstring)

	if err != nil {
		panic(err)
	}
	defer session.Close()

	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var p Postload
		var d Webhook
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
		} else { //Create data in database:
			d.BaseCurrency = p.BaseCurrency
			d.MaxTriggerValue = p.MaxTriggerValue
			d.MinTriggerValue = p.MinTriggerValue
			d.TargetCurrency = p.TargetCurrency
			d.WebhookURL = p.WebhookURL
			d.CurrentRate = 0 //Set equal to value in database relative to base and TargetCurrency
			err := session.DB(tempstring).C("testcollection").Insert(d)
			if err != nil {
				fmt.Fprintln(w, "Error in Insert()", err.Error())
			}
			id := session.DB(tempstring).C("testcollection").Find(bson.M{"targetCurrency": "NOK"})
			fmt.Fprintln(w, id)
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