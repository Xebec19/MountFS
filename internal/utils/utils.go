package utils

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/Xebec19/reimagined-lamp/pkg/logger"
)

var once sync.Once
var appDir string

func GetAppDir() string {
	once.Do(func() {
		rootDir, err := os.UserConfigDir()
		if err != nil {
			logger.Error("Failed to get user config directory: ", err)
			return
		}

		appDir = filepath.Join(rootDir, "mountfs")
		if err := os.MkdirAll(appDir, 0700); err != nil {
			logger.Error("Failed to create app directory: ", err)
			return
		}

		logger.Info("Root dir set ", appDir)
	})
	return appDir
}

func GetConfigPath() string {
	GetAppDir()

	return filepath.Join(appDir, "config.json")
}

func GetTokenPath() string {
	if dir := GetAppDir(); dir == "" {
		return ""
	}

	return filepath.Join(appDir, "token.json")
}
