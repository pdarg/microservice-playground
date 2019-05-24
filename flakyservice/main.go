package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	HTTPListenPort = 8080
)

var beFlaky bool

func main() {

	beFlaky = false

	router := mux.NewRouter()
	router.HandleFunc("/", Index)
	router.HandleFunc("/start", Start)
	router.HandleFunc("/stop", Stop)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", HTTPListenPort), router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	if beFlaky {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Oops")
		return
	}
	log.Println("Success")
	fmt.Fprintln(w, "SUCCESS")
}

func Start(w http.ResponseWriter, r *http.Request) {
	beFlaky = false
	fmt.Fprintln(w, "Flaky Service is now working")
	log.Println("Back on track!")
}

func Stop(w http.ResponseWriter, r *http.Request) {
	beFlaky = true
	fmt.Fprintln(w, "Flaky Service is now flaky")
	log.Println("Look a squirrel!")
}
