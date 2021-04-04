package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

var lineSpeedAggr LineSpeedAggr

func main() {
	lineSpeedAggr = make(LineSpeedAggr)

	r := mux.NewRouter()

	r.HandleFunc("/linespeed", linespeed).Methods("POST")
	r.HandleFunc("/metrics", metrics).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./dist/")))

	http.Handle("/", r)

	_ = http.ListenAndServe(":843", nil)
}

func linespeed(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var ls LineSpeed
	err := decoder.Decode(&ls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lineSpeedAggr.AddLs(ls)

	resStatus := http.StatusCreated

	if ls.IsOld() {
		resStatus = http.StatusPartialContent
	}

	w.WriteHeader(resStatus)
}

func metrics(w http.ResponseWriter, _ *http.Request) {
	js, err := json.Marshal(lineSpeedAggr.Metrics())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}
