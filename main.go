package main

import (
	"fmt"
	"net/http"
	"os"
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
	switch r.Method {
	case "POST":

	case "GET":

	case "DELETE":

	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

	fmt.Fprintln(w, "hello, world")

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
