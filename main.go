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
	fmt.Fprintln(w, "hello, world")

	timer := time.NewTimer(time.Second * 10).C

	timerFinished := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		timerFinished <- true
	}()

	for {
		select {
		case <-timer:
			fmt.Println("Timer expired")
		case <-timerFinished:
			fmt.Println("Done")
		}
	}

}
