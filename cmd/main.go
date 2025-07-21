package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func SaveToken(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	code := params.Get("code")

	w.Write([]byte(code))
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", SaveToken)

	log.Fatal(http.ListenAndServe(":8080", r))
}
