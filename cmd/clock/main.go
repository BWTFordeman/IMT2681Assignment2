package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"
)

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

	fmt.Println(f.BaseCurrency, f.Date, f.Rates)
	//k := time.Now().AddDate(0, 0, -3)

	//Connect to database:
	USER := os.Getenv("DB_USER")
	PASSWORD := os.Getenv("DB_PASSWORD")
	//DBNAME := os.Getenv("DB_NAME")
	tempstring := ("mongodb://" + USER + ":" + PASSWORD + "@ds241055.mlab.com:41055/imt2681")

	session, err := mgo.Dial(tempstring)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	//add part to databse:
	//if date of object is less than k then remove from database
}

func main() {

	getFixerData()
	start := time.Now()
	for {
		//time.Now().String()
		delay := time.Minute * 20
		elapsed := time.Now().Sub(start)
		fmt.Println("\n", time.Now(), "\n", elapsed)
		if elapsed > time.Hour*24 {
			start = time.Now()
			getFixerData()
			//Check if data gotten from fixer exceedes someones threshold in webhook table.
		}
		time.Sleep(delay)
	}
}
