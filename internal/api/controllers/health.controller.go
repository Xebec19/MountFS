package controllers

import "net/http"

func TestHealth(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("alive"))
}
