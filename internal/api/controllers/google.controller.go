package controllers

import (
	"net/http"
	"os"

	"github.com/Xebec19/reimagined-lamp/internal/utils"
	"github.com/Xebec19/reimagined-lamp/pkg/logger"
)

func SaveToken(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	code := params.Get("code")

	tokenFilePath := utils.GetTokenPath()

	if err := os.WriteFile(tokenFilePath, []byte(code), 0600); err != nil {
		logger.Error("Failed to write token file", "error", err, "path", tokenFilePath)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info("code written in ", tokenFilePath)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(code))
}
