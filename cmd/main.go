package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/Xebec19/reimagined-lamp/internal/utils"
	"github.com/gorilla/mux"
)

func SaveToken(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	code := params.Get("code")

	tokenFilePath := utils.GetTokenPath()

	os.WriteFile(tokenFilePath, []byte(code), 0700)

	slog.Info("code written in ", tokenFilePath)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(code))
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", SaveToken)

	log.Fatal(http.ListenAndServe("127.0.0.1"+":8080", r))
}
