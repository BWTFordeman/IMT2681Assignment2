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
type Fixer struct {
	BaseCurrency string             `json:"base"`
	Date         string             `json:"date"`
	Rates        map[string]float64 `json:"rates"`
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
	k := time.Now()
	k.AddDate(0, 0, -3)
	err = session.DB(DBNAME).C("fixerdata").Remove(bson.M{"date": k.Format("2006-01-02")})
	if err != nil {
		fmt.Println("fixerdata - No data older than 3 days")
	}
	//If date doesn't already exist - add part to database:
	d := Fixer{}
	t := time.Now()

	err = session.DB(DBNAME).C("fixerdata").Find(bson.M{"date": t.Format("2006-01-02")}).One(&d)
	if err != nil {
		fmt.Println("date", f.Date)
		err = session.DB(DBNAME).C("fixerdata").Insert(bson.M{"baseCurrency": f.BaseCurrency, "date": f.Date, "rates": f.Rates})
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

		//Goes through all targetcurrencies and checks if currentRate is same(prints out they are same)
		//if not same{put the f.Rates value into database}
		fmt.Println("Here comes the updating:")
		for _, k := range web {
			for j, l := range f.Rates {
				if k.TargetCurrency == j {
					fmt.Println("Updating", k.WebhookURL)
					update := bson.M{"$set": bson.M{"currentRate": l}}
					err := session.DB(DBNAME).C("webhooks").UpdateId(k.ID, update)
					if err != nil {
						fmt.Println("Could not update webhook")
					}
				}
			}
		}
	}
}

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
		//Check the values and invokeWebhook when required
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
