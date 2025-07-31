package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Xebec19/reimagined-lamp/internal/api/server"
	"github.com/Xebec19/reimagined-lamp/internal/utils"
	"github.com/Xebec19/reimagined-lamp/pkg/logger"
)

func main() {

	portArg := os.Getenv("PORT")

	port, err := utils.ConvertFromStrToUint(portArg)
	if err != nil {
		logger.Error("Could not get PORT")
		os.Exit(1)
	}
	srv, err := server.CreateServer(uint(port))
	if err != nil {
		logger.Error("Could not create server ", err)
	}

	go func() {
		srv.StartServer()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	err = srv.ShutdownServer()
	if err != nil {
		logger.Error("Could not shut down server ", err)
	}

}
