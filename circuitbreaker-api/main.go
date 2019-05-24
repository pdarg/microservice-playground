package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sony/gobreaker"
)

const (
	HTTPListenPort     = 8081
	CBTimeoutSeconds   = 15
	CBMinRequests      = 3
	CBFailureThreshold = 0.6
	BackendServiceURL  = "http://localhost:8080"
)

type Result struct {
	Message string
}

var cb *gobreaker.CircuitBreaker

func initCircuitBreaker() {
	var st gobreaker.Settings
	st.Name = "Flaky Service GET"
	st.Timeout = time.Duration(CBTimeoutSeconds) * time.Second
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= CBMinRequests && failureRatio >= CBFailureThreshold
	}

	cb = gobreaker.NewCircuitBreaker(st)
}

func Get(url string) (*http.Response, error) {
	response, err := cb.Execute(func() (interface{}, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, errors.New("Backend failed")
		}

		return resp, nil
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("Success")
	return response.(*http.Response), nil
}

func main() {

	initCircuitBreaker()

	router := mux.NewRouter()
	router.HandleFunc("/", Index)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", HTTPListenPort), router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	result := Result{}
	response, err := Get(BackendServiceURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result.Message = fmt.Sprintf("Failed to connect to flakyservice: %s", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			panic(err)
		}
		return
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		result.Message = fmt.Sprintf("Failed to parse response from flakyservice: %s", err)
	} else if response.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
		result.Message = fmt.Sprintf("Failed get valid response from flakyservice: Got %d response", response.StatusCode)
	} else if len(contents) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		result.Message = "Failed get valid response from flakyservice: No data"
	} else {
		result.Message = strings.TrimSpace(string(contents))
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}
