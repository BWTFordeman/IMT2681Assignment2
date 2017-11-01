package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

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
		fmt.Fprintln(w, d)
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
			} else {
				id := bson.NewObjectId() //Get currentRate from fixerdata collection and put that in currentRate instead of 0
				err := session.DB(DBNAME).C("webhooks").Insert(bson.M{"_id": id, "webhookURL": p.WebhookURL, "baseCurrency": p.BaseCurrency, "targetCurrency": p.TargetCurrency, "maxTriggerValue": p.MaxTriggerValue, "minTriggerValue": p.MinTriggerValue, "currentRate": 0})
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
