package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		fmt.Println(err.Error(), " Panic or something")
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	//if r.Method == "POST" {}
	//else if r.Method == "GET" {}
	//else if r.Method == "DELETE"{}
	//else {http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)}
	fmt.Fprintln(w, "hello, world")

	startTime := time.Now()

	t := time.Now()
	elapsed := t.Sub(startTime)
	if elapsed < time.Second*10 {
		fmt.Fprintln(w, "ok", time.Now())
	}

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
