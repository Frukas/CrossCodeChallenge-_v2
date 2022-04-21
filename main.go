package main

import (
	"encoding/json"
	"log"

	//"log"
	"time"

	"net/http"

	dataretriver "github.com/frukas/crosscodechallenge/DataRetriver"
	mergesort "github.com/frukas/crosscodechallenge/MergeSort"
	"github.com/gorilla/mux"
)

type restults struct {
	Numbers       []float32 `json:"numbers"`
	Pages         int       `json:"pages"`
	ElapseMinutes float64   `json:"elapseminutes"`
}

func NumbersResult(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	testResult, pages := dataretriver.PageDataRetrive()

	mergesort.MergeSort(testResult)

	elapsed := time.Since(start)

	response := restults{
		Numbers:       testResult,
		Pages:         pages,
		ElapseMinutes: elapsed.Minutes(),
	}

	json.NewEncoder(w).Encode(response)
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/NumbersResult", NumbersResult)

	log.Fatal(http.ListenAndServe(":8080", router))
}
