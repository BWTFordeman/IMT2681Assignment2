package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
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
	ID              bson.ObjectId `json:"_id" bson:"_id"`
	WebhookURL      string        `json:"webhookURL" bson:"webhookURL"`
	BaseCurrency    string        `json:"baseCurrency" bson:"baseCurrency"`
	TargetCurrency  string        `json:"targetCurrency" bson:"targetCurrency"`
	MinTriggerValue float32       `json:"minTriggerValue" bson:"minTriggerValue"`
	MaxTriggerValue float32       `json:"maxTriggerValue" bson:"maxTriggerValue"`
	CurrentRate     float32       `json:"currentRate" bson:"currentRate"`
}

/*Fixer retrieves data from fixer collection:				//Rates map[string]float64 `json:"rates"`
Use for i := range rates {

}*/
type Fixer struct {
	BaseCurrency string `json:"base"`
	Date         string `json:"date"`
	Rates        struct {
		AUD float32 `json:"AUD"`
		BGN float32 `json:"BGN"`
		BRL float32 `json:"BRL"`
		CAD float32 `json:"CAD"`
		CHF float32 `json:"CHF"`
		CNY float32 `json:"CNY"`
		CZK float32 `json:"CZK"`
		DKK float32 `json:"DKK"`
		GBP float32 `json:"GBP"`
		HRK float32 `json:"HRK"`
		HKD float32 `json:"HKD"`
		HUF float32 `json:"HUF"`
		IDR float32 `json:"IDR"`
		ILS float32 `json:"ILS"`
		INR float32 `json:"INR"`
		JPY float32 `json:"JPY"`
		KRW float32 `json:"KRW"`
		MXN float32 `json:"MXN"`
		MYR float32 `json:"MYR"`
		NOK float32 `json:"NOK"`
		NZD float32 `json:"NZD"`
		PHP float32 `json:"PHP"`
		PLN float32 `json:"PLN"`
		RON float32 `json:"RON"`
		RUB float32 `json:"RUB"`
		SEK float32 `json:"SEK"`
		SGD float32 `json:"SGD"`
		THB float32 `json:"THB"`
		TRY float32 `json:"TRY"`
		USD float32 `json:"USD"`
		ZAR float32 `json:"ZAR"`
	} `json:"rates"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{id}", getWebhooks).Methods("GET")
	r.HandleFunc("/{id}", deleteWebhooks).Methods("DELETE")
	http.Handle("/", r)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		fmt.Println(err.Error(), "Panic or something")
	}
}

func deleteWebhooks(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	url2 := strings.Split(url, "/")
	//Check that url2[1] = 24 values long for id

	USER := os.Getenv("DB_USER")
	PASSWORD := os.Getenv("DB_PASSWORD")
	DBNAME := os.Getenv("DB_NAME")
	tempstring := ("mongodb://" + USER + ":" + PASSWORD + "@ds241055.mlab.com:41055/imt2681")

	session, err := mgo.Dial(tempstring)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.DB(DBNAME).C("webhooks").Remove(bson.M{"_id": bson.ObjectIdHex(url2[1])})
	if err != nil {
		http.Error(w, "Could not find any object with that id", http.StatusBadRequest)
	}
}

func getWebhooks(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	url2 := strings.Split(url, "/")
	//Check that url2[1] = 24 values long for id

	USER := os.Getenv("DB_USER")
	PASSWORD := os.Getenv("DB_PASSWORD")
	DBNAME := os.Getenv("DB_NAME")
	tempstring := ("mongodb://" + USER + ":" + PASSWORD + "@ds241055.mlab.com:41055/imt2681")

	session, err := mgo.Dial(tempstring)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	d := Webhook{}
	err = session.DB(DBNAME).C("webhooks").Find(bson.M{"_id": bson.ObjectIdHex(url2[1])}).One(&d)
	if err != nil {
		http.Error(w, "Object doesn't exist", http.StatusBadRequest)
	} else {
		fmt.Fprintln(w, "{\n\tbaseCurrency", d.BaseCurrency, "\n\ttargetCurrency:", d.TargetCurrency, "\n\tcurrentRate:", d.CurrentRate, "\n\tminTriggerValue:", d.MinTriggerValue, "\n\tmaxTriggerValue:", d.MaxTriggerValue, "\n}")
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	lang := [...]string{"EUR", "AUD", "BGN", "BRL", "CAD", "CHF", "CNY", "CZK", "DKK", "GBP", "HKD", "HRK", "HUF", "IDR", "ILS", "INR", "JPY", "KRW",
		"MXN", "MYR", "NOK", "NZD", "PHP", "PLN", "RON", "RUB", "SEK", "SGD", "THB", "TRY", "USD", "ZAR"}
	USER := os.Getenv("DB_USER")
	PASSWORD := os.Getenv("DB_PASSWORD")
	DBNAME := os.Getenv("DB_NAME")
	tempstring := ("mongodb://" + USER + ":" + PASSWORD + "@ds241055.mlab.com:41055/imt2681")

	session, err := mgo.Dial(tempstring)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var p Postload

		err := decoder.Decode(&p)
		if err != nil {
			fmt.Fprintln(w, "Error decoding webhook post", err.Error())
		}

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

		//Create object in database if valid:
		if err != nil || base != true || target != true {
			http.Error(w, "Invalid post value", http.StatusBadRequest)
		} else {
			//Create data in database if not there from before:
			d := Webhook{}
			err = session.DB(DBNAME).C("webhooks").Find(bson.M{"webhookURL": p.WebhookURL, "targetCurrency": p.TargetCurrency}).One(&d)
			if err == nil {
				http.Error(w, "Object already exists", http.StatusBadRequest)
			} else { //Get currentRate from fixerdata collection and put that in currentRate.
				f := Fixer{}
				k := time.Now()
				err = session.DB(DBNAME).C("fixerdata").Find(bson.M{"date": k.String()}).One(&f)
				if err != nil {
					fmt.Fprintln(w, "Could not get currentRate data")
				}
				id := bson.NewObjectId()
				fmt.Fprintln(w, "current value", getCurrentValue(f, p.TargetCurrency))
				err := session.DB(DBNAME).C("webhooks").Insert(bson.M{"_id": id, "webhookURL": p.WebhookURL, "baseCurrency": p.BaseCurrency, "targetCurrency": p.TargetCurrency, "maxTriggerValue": p.MaxTriggerValue, "minTriggerValue": p.MinTriggerValue, "currentRate": getCurrentValue(f, p.TargetCurrency)})
				if err != nil {
					fmt.Fprintln(w, "Error in Insert()", err.Error())
				}

				d = Webhook{}
				err = session.DB(DBNAME).C("webhooks").Find(bson.M{"webhookURL": p.WebhookURL, "targetCurrency": p.TargetCurrency}).One(&d)
				fmt.Fprintln(w, "Id of your webhook:", d.ID.Hex())
			}
		}

		defer r.Body.Close()
	} else { //If not post:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

//THIS NEEDS TO BE FIXED AS FAST AS POSSIBLE!! NOT GOOD!
func getCurrentValue(f Fixer, targetCurrency string) float32 {
	if targetCurrency == "AUD" {
		return f.Rates.AUD
	} else if targetCurrency == "BGN" {
		return f.Rates.BGN
	} else if targetCurrency == "BRL" {
		return f.Rates.BRL
	} else if targetCurrency == "CAD" {
		return f.Rates.CAD
	} else if targetCurrency == "CHF" {
		return f.Rates.CHF
	} else if targetCurrency == "CNY" {
		return f.Rates.CNY
	} else if targetCurrency == "CZK" {
		return f.Rates.CZK
	} else if targetCurrency == "DKK" {
		return f.Rates.DKK
	} else if targetCurrency == "GBP" {
		return f.Rates.GBP
	} else if targetCurrency == "HKD" {
		return f.Rates.HKD
	} else if targetCurrency == "HRK" {
		return f.Rates.HRK
	} else if targetCurrency == "HUF" {
		return f.Rates.HUF
	} else if targetCurrency == "IDR" {
		return f.Rates.IDR
	} else if targetCurrency == "ILS" {
		return f.Rates.ILS
	} else if targetCurrency == "INR" {
		return f.Rates.INR
	} else if targetCurrency == "JPY" {
		return f.Rates.JPY
	} else if targetCurrency == "KRW" {
		return f.Rates.KRW
	} else if targetCurrency == "MXN" {
		return f.Rates.MXN
	} else if targetCurrency == "MYR" {
		return f.Rates.MYR
	} else if targetCurrency == "NOK" {
		return f.Rates.NOK
	} else if targetCurrency == "NZD" {
		return f.Rates.NZD
	} else if targetCurrency == "PHP" {
		return f.Rates.PHP
	} else if targetCurrency == "PLN" {
		return f.Rates.PLN
	} else if targetCurrency == "RON" {
		return f.Rates.RON
	} else if targetCurrency == "RUB" {
		return f.Rates.RUB
	} else if targetCurrency == "SEK" {
		return f.Rates.SEK
	} else if targetCurrency == "SGD" {
		return f.Rates.SGD
	} else if targetCurrency == "THB" {
		return f.Rates.THB
	} else if targetCurrency == "TRY" {
		return f.Rates.TRY
	} else if targetCurrency == "USD" {
		return f.Rates.USD
	} else if targetCurrency == "ZAR" {
		return f.Rates.ZAR
	}
	return 0
}
