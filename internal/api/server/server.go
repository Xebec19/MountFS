package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Xebec19/reimagined-lamp/internal/api/routes"
	"github.com/Xebec19/reimagined-lamp/pkg/logger"
)

type ApiServer struct {
	srv *http.Server
}

func (s ApiServer) StartServer() {

	go func() {
		logger.Info("Starting server on port", s.srv.Addr)

		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server")
}

func (s ApiServer) ShutdownServer() error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	return s.srv.Shutdown(ctx)
}

func CreateServer(port uint) (Server, error) {

	r := routes.NewRouter()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	return ApiServer{srv: srv}, nil
}
