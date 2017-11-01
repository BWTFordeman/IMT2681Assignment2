package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Webhook retrieves data from the webhook collection:
type Webhook struct {
	ID              bson.ObjectId `json:"_id" bson:"_id"`
	WebhookURL      string        `json:"webhookURL"`
	BaseCurrency    string        `json:"baseCurrency"`
	TargetCurrency  string        `json:"targetCurrency"`
	MinTriggerValue float32       `json:"minTriggerValue"`
	MaxTriggerValue float32       `json:"maxTriggerValue"`
	CurrentRate     float32       `json:"currentRate"`
}

//Fixer retrieves latest data from fixer.io
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

//getFixerData retrieves data from fixer.io and puts them in collection "fixerdata"
func getFixerData() {
	var f Fixer
	url := "http://api.fixer.io/latest?base=EUR"

	link, err := http.Get(url)
	if err != nil {
		fmt.Print("1", err.Error())
	}

	defer link.Body.Close()

	body, err := ioutil.ReadAll(link.Body)
	if err != nil {
		fmt.Print("2", err.Error())
	}

	err = json.Unmarshal(body, &f)
	if err != nil {
		fmt.Print("3", err.Error())
	}
	k := time.Now().AddDate(0, 0, -3)

	//Connect to database:
	USER := os.Getenv("DB_USER")
	PASSWORD := os.Getenv("DB_PASSWORD")
	DBNAME := os.Getenv("DB_NAME")
	tempstring := ("mongodb://" + USER + ":" + PASSWORD + "@ds241055.mlab.com:41055/imt2681")
	session, err := mgo.Dial(tempstring)
	if err != nil {
		fmt.Println("Error connecting to database", err.Error())
	}
	defer session.Close()

	//If date doesn't already exist - add part to database:
	d := Fixer{}
	err = session.DB(DBNAME).C("fixerdata").Find(bson.M{"date": f.Date}).One(&d)
	if err != nil {
		err = session.DB(DBNAME).C("fixerdata").Insert(bson.M{"baseCurrency": f.BaseCurrency, "date": f.Date, "rates": f.Rates})
		if err != nil {
			fmt.Println("Error in Insert()", err.Error())
		}
	} else {
		fmt.Println("fixerdata - Found object in database")
	}
	//if date of object is less than k then remove from database
	err = session.DB(DBNAME).C("fixerdata").Remove(bson.M{"date": bson.M{"$lt": k.String()}})
	if err != nil {
		fmt.Println("fixerdata - No data older than 3 days")
	}
}

//Check if data gotten from fixer exceedes someones threshold in webhook table.
func updateWebhooks() {

	//Connect to database:
	USER := os.Getenv("DB_USER")
	PASSWORD := os.Getenv("DB_PASSWORD")
	DBNAME := os.Getenv("DB_NAME")
	tempstring := ("mongodb://" + USER + ":" + PASSWORD + "@ds241055.mlab.com:41055/imt2681")
	session, err := mgo.Dial(tempstring)
	if err != nil {
		fmt.Println("Error connecting to database", err.Error())
	}
	defer session.Close()

	//Go through all webhooks
	d := []Webhook{}
	err = session.DB(DBNAME).C("webhooks").Find(nil).All(&d)
	if err != nil {
		//some error
	} else {
		fmt.Println("webhooks - found these:", d[0])
	}
}

func main() {

	getFixerData()
	updateWebhooks()
	start := time.Now()
	for {
		//time.Now().String()
		delay := time.Minute * 20
		elapsed := time.Now().Sub(start)
		if elapsed > time.Hour*24 {
			start = time.Now()
			getFixerData()
			updateWebhooks()
		}
		time.Sleep(delay)
	}
}
