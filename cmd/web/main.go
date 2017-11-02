package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

//Latest gets the latest value
type Latest struct {
	BaseCurrency   string `json:"baseCurrency"`
	TargetCurrency string `json:"targetCurrency"`
}

//Fixer retrieves data from fixer collection:
type Fixer struct {
	BaseCurrency string             `json:"base"`
	Date         string             `json:"date"`
	Rates        map[string]float64 `json:"rates"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", root)
	r.HandleFunc("/{id}", getWebhooks).Methods("GET")
	r.HandleFunc("/{id}", deleteWebhooks).Methods("DELETE")
	r.HandleFunc("/latest", getLatest).Methods("POST")
	r.HandleFunc("/average", getAverage).Methods("POST")
	r.HandleFunc("/evaluationtrigger", triggerwebhooks).Methods("GET")
	http.Handle("/", r)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		fmt.Println(err.Error(), "Panic or something")
	}
}

func getAverage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var l Latest

	err := decoder.Decode(&l)
	if err != nil {
		http.Error(w, "Error decoding post request for average", http.StatusBadRequest)
	} else {
		//Connecting to database:
		USER := os.Getenv("DB_USER")
		PASSWORD := os.Getenv("DB_PASSWORD")
		DBNAME := os.Getenv("DB_NAME")
		tempstring := ("mongodb://" + USER + ":" + PASSWORD + "@ds241055.mlab.com:41055/imt2681")

		session, err := mgo.Dial(tempstring)
		if err != nil {
			panic(err)
		}
		defer session.Close()
		//Get values of 3 days and sends average:
		f := []Fixer{}

		err = session.DB(DBNAME).C("fixerdata").Find(nil).All(&f)
		if err != nil {
			http.Error(w, "Could not get average value", http.StatusBadRequest)
		} else {
			var total float64
			var amount float64
			total = 0
			amount = 0
			for _, k := range f { //range over 3 days
				for y, j := range k.Rates { //Range over all languages
					if y == l.TargetCurrency {
						total = total + j
						amount++
					}
				}
			}
			fmt.Fprintln(w, total/amount)
		}
	}
}

//triggerwebhooks sends messages to all webhooks that have current value breaking the threshold
func triggerwebhooks(w http.ResponseWriter, r *http.Request) {
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

//invokeWebhook sends messages for one webhook in the system
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

func getLatest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var l Latest

	err := decoder.Decode(&l)
	if err != nil {
		http.Error(w, "Error decoding post request for latest", http.StatusBadRequest)
	} else {
		//Connecting to database:
		USER := os.Getenv("DB_USER")
		PASSWORD := os.Getenv("DB_PASSWORD")
		DBNAME := os.Getenv("DB_NAME")
		tempstring := ("mongodb://" + USER + ":" + PASSWORD + "@ds241055.mlab.com:41055/imt2681")

		session, err := mgo.Dial(tempstring)
		if err != nil {
			panic(err)
		}
		defer session.Close()

		//Getting current value:
		f := Fixer{}
		t := time.Now()
		if t.Hour() < 17 {
			t = t.AddDate(0, 0, -1)
		}
		err = session.DB(DBNAME).C("fixerdata").Find(bson.M{"date": t.Format("2006-01-02")}).One(&f)
		if err != nil {
			fmt.Fprintln(w, "Could not get currentRate data")
		}

		fmt.Fprintln(w, getCurrentValue(f, l.TargetCurrency))
	}
}

func deleteWebhooks(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	url2 := strings.Split(url, "/")
	//TODO Check that url2[1] = 24 values long for id

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
	} else {
		http.Error(w, "Deleted object", http.StatusAccepted)
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

				if k.Hour() < 17 {
					k = k.AddDate(0, 0, -1)
				}

				err = session.DB(DBNAME).C("fixerdata").Find(bson.M{"date": k.Format("2006-01-02")}).One(&f)
				if err != nil {
					//fmt.Println("Currentdata set to 0, could not find current value")
				}
				id := bson.NewObjectId()
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

//Returns float value of current data
func getCurrentValue(f Fixer, targetCurrency string) float64 {

	for i, k := range f.Rates {
		if targetCurrency == i {
			return k
		}
	}
	return 0
}
