package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

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

//Fixer retrieves latest data from fixer.io
//Rates should DEFINITELY BE DONE DIFFERENTLY!! BUT NOT ENOUGH TIME TO DO IT FOR ASSIGNMENT2!! TOO MUCH code for such little work...Hard to think of a solution with this stress
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

//getFixerData retrieves data from fixer.io and puts them in collection "fixerdata"
func getFixerData() Fixer {
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

	//if date of object is less than k then remove from database
	err = session.DB(DBNAME).C("fixerdata").Remove(bson.M{"date": bson.M{"$lt": k.String()}})
	if err != nil {
		fmt.Println("fixerdata - No data older than 3 days")
	}

	//If date doesn't already exist - add part to database:
	d := Fixer{}
	err = session.DB(DBNAME).C("fixerdata").Find(bson.M{"date": time.Now()}).One(&d)
	if err != nil {
		err = session.DB(DBNAME).C("fixerdata").Insert(bson.M{"baseCurrency": f.BaseCurrency, "date": time.Now(), "rates": f.Rates})
		if err != nil {
			fmt.Println("Error in Insert()", err.Error())
		}
		return f
	}
	fmt.Println("fixerdata - Found object in database")
	return d
}

//Check if data gotten from fixer exceedes someones threshold in webhook table.
func updateWebhooks(f Fixer) {

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
	web := []Webhook{}
	err = session.DB(DBNAME).C("webhooks").Find(nil).All(&web)
	if err != nil {
		fmt.Println("webhooks - Could not find any webhooks")
	} else {
		//d[0].CurrentRate check if this is equal to the found one. if not put in right value
		for i := range web {
			if web[i].TargetCurrency == "AUD" {
				if web[i].CurrentRate == f.Rates.AUD {
					fmt.Println("currentRate is of same value")
				} else { //update database
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.AUD}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "BGN" {
				if web[i].CurrentRate == f.Rates.BGN {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.BGN}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "BRL" {
				if web[i].CurrentRate == f.Rates.BRL {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.BRL}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "CAD" {
				if web[i].CurrentRate == f.Rates.CAD {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.CAD}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "CHF" {
				if web[i].CurrentRate == f.Rates.CHF {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.CHF}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "CNY" {
				if web[i].CurrentRate == f.Rates.CNY {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.CNY}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "CZK" {
				if web[i].CurrentRate == f.Rates.CZK {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.CZK}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "DKK" {
				if web[i].CurrentRate == f.Rates.DKK {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.DKK}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "GBP" {
				if web[i].CurrentRate == f.Rates.GBP {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.GBP}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "HKD" {
				if web[i].CurrentRate == f.Rates.HKD {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.HKD}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "HRK" {
				if web[i].CurrentRate == f.Rates.HRK {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.HRK}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "HUF" {
				if web[i].CurrentRate == f.Rates.HUF {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.HUF}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "IDR" {
				if web[i].CurrentRate == f.Rates.IDR {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.IDR}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "ILS" {
				if web[i].CurrentRate == f.Rates.ILS {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.ILS}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "INR" {
				if web[i].CurrentRate == f.Rates.INR {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.INR}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "JPY" {
				if web[i].CurrentRate == f.Rates.JPY {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.JPY}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "KRW" {
				if web[i].CurrentRate == f.Rates.KRW {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.KRW}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "MXN" {
				if web[i].CurrentRate == f.Rates.MXN {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.MXN}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "MYR" {
				if web[i].CurrentRate == f.Rates.MYR {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.MYR}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "NOK" {
				if web[i].CurrentRate == f.Rates.NOK {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.NOK}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "NZD" {
				if web[i].CurrentRate == f.Rates.NZD {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.NZD}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "PHP" {
				if web[i].CurrentRate == f.Rates.PHP {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.PHP}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "PLN" {
				if web[i].CurrentRate == f.Rates.PLN {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.PLN}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "RON" {
				if web[i].CurrentRate == f.Rates.RON {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.RON}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "RUB" {
				if web[i].CurrentRate == f.Rates.RUB {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.RUB}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "SEK" {
				if web[i].CurrentRate == f.Rates.SEK {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.SEK}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "SGD" {
				if web[i].CurrentRate == f.Rates.SGD {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.SGD}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "THB" {
				if web[i].CurrentRate == f.Rates.THB {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.THB}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "TRY" {
				if web[i].CurrentRate == f.Rates.TRY {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.TRY}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "USD" {
				if web[i].CurrentRate == f.Rates.USD {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.USD}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			} else if web[i].TargetCurrency == "ZAR" {
				if web[i].CurrentRate == f.Rates.ZAR {
					fmt.Println("currentRate is of same value")
				} else {
					update := bson.M{"$set": bson.M{"currentRate": f.Rates.ZAR}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(web[i].ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			}
		}
	}
}

//Check if d[0].CurrentRate extends the threshold
//	then send through webhook
//Checks if any current values are beyond anyones threshold and thereafter send messages.
func sendToWebhooks() {
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

	web := []Webhook{}
	err = session.DB(DBNAME).C("webhooks").Find(nil).All(&web)
	if err != nil {
		fmt.Println("webhooks - Could not find any webhooks")
	} else {

		for i := range web {
			if web[i].CurrentRate > web[i].MaxTriggerValue || web[i].CurrentRate < web[i].MinTriggerValue {

				invokeWebhook(web[i].WebhookURL, web[i].TargetCurrency, web[i].CurrentRate, web[i].MinTriggerValue, web[i].MaxTriggerValue)
			}
		}
	}
}

//TODO add validation for /slack or /github at end of webhookURL
//invokeWebhook sends messages through webhooks created in the system
func invokeWebhook(webhookURL string, targetCurrency string, currentRate float32, minTriggerValue float32, maxTriggerValue float32) {
	current := strconv.FormatFloat(float64(currentRate), 'f', 2, 32)
	mintrigger := strconv.FormatFloat(float64(minTriggerValue), 'f', 2, 32)
	maxtrigger := strconv.FormatFloat(float64(maxTriggerValue), 'f', 2, 32)
	res, err := http.PostForm(webhookURL, url.Values{"content": {"{\n\tbaseCurrency: EUR" + "\n\ttargetCurrency:\t" + targetCurrency + "\n\tcurrentRate:\t" + current + "\n\tminTriggerValue:\t" + mintrigger + "\n\tmaxTriggerValue:\t" + maxtrigger + "\n}"}, "username": {"CurrencyChecker"}})
	if err != nil {
		fmt.Println("Error posting webhook message")
	} else {
		fmt.Println("A webhook message is sent")
	}
	if res.StatusCode == 200 || res.StatusCode == 204 {
		fmt.Println("statuscode: ", res.StatusCode)
	} else {
		fmt.Println("Wrong status: ", res.StatusCode, http.StatusText(res.StatusCode))
	}
}

func main() {
	delay := time.Minute * 20
	f := getFixerData()
	updateWebhooks(f)
	sendToWebhooks()

	//Timer:
	for {
		time.Sleep(delay)
		g := getFixerData()
		updateWebhooks(g)
		sendToWebhooks()
	}
}
