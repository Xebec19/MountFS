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

	// Validate the code parameter
	if code == "" {
		logger.Error("Missing or empty code parameter")
		http.Error(w, "Bad request: missing code parameter", http.StatusBadRequest)
		return
	}

	// Basic validation for OAuth authorization code format
	if len(code) < 10 || len(code) > 512 {
		logger.Error("Invalid code parameter length", "length", len(code))
		http.Error(w, "Bad request: invalid code format", http.StatusBadRequest)
		return
	}

	tokenFilePath := utils.GetTokenPath()

	if err := os.WriteFile(tokenFilePath, []byte(code), 0600); err != nil {
		logger.Error("Failed to write token file", "error", err, "path", tokenFilePath)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Info("Token saved successfully", "path", tokenFilePath)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Token saved successfully"))
}
