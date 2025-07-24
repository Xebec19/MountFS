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
		rootDir, _ := os.UserConfigDir()
		appDir = filepath.Join(rootDir, "mountfs")
		os.MkdirAll(appDir, 0700)
		logger.Info("Root dir set ", appDir)
	})

	return appDir
}

func GetConfigPath() string {
	GetAppDir()

	return filepath.Join(appDir, "config.json")
}

func GetTokenPath() string {
	GetAppDir()

	return filepath.Join(appDir, "token.json")
}
