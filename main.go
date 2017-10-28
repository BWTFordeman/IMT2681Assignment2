package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

//Postload datatype I get from fixer.io
type Postload struct {
	WebhookURL      string `json:"webhookURL"`
	BaseCurrency    string `json:"baseCurrency"`
	TargetCurrency  string `json:"targetCurrency"`
	MinTriggerValue int    `json:"minTriggerValue"`
	MaxTriggerValue int    `json:"maxTriggerValue"`
}

/*
//Make another one for requesting data from fixer.io
type Payload struct {
	WebhookURL   string `json:"url"`
	BaseCurrency string `json:"baseCurrency"`
	Date         int    `json:"date"`
	Rates        struct {
		currency map[string]int
	} `json:"rates"`
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
	switch r.Method {
	case "POST":
		decoder := json.NewDecoder(r.Body)
		var p Postload
		err := decoder.Decode(&p)
		fmt.Fprintln(w, err.Error())
		if err != nil {
			http.Error(w, "Invalid post value", http.StatusBadRequest)
		}
		defer r.Body.Close()
		fmt.Fprintln(w, p.WebhookURL, p.BaseCurrency, p.TargetCurrency)
		//Can do stuff with p. ...

	case "GET":
		fmt.Fprintln(w, "GET")

	case "DELETE":

	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	//Do stuff every 24 hour:

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
