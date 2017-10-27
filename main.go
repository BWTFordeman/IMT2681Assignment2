package main

import (
	"log"
	"net/http"
	"os"
)

//Payload datatype I get from fixer.io
/*type Payload struct {
	base  string `json:"base"`
	date  int    `json:"date"`
	rates struct {
		currency map[string]int
	} `json:"rates"`
}*/

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Println("panic - ListenAndServe(14)")
	}
	//log.Println("http.ListenAndServe", http.ListenAndServe(":"+os.Getenv("PORT"), nil), nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

}
