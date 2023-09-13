package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"signature-server/logger"
	"time"
)

// Server ...
type Server struct {
	name    string
	port    int
	timeout time.Duration
	cleanup func()
	handler http.Handler
}

// NewServer ...
func NewServer(name string, port int, timeout time.Duration, h http.Handler) *Server {
	return &Server{
		name:    name,
		port:    port,
		timeout: timeout,
		handler: h,
	}
}

// Run ...
func (svr *Server) Run() {
	logger.Infof("starting %s server...", svr.name)

	server := http.Server{
		Addr:              fmt.Sprintf(":%d", svr.port),
		Handler:           svr.handler,
		ReadTimeout:       svr.timeout,
		ReadHeaderTimeout: svr.timeout,
		WriteTimeout:      svr.timeout,
		IdleTimeout:       svr.timeout,
	}

	go func() {
		logger.Infof("%s server listening on port %d", svr.name, svr.port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("%s server stopped listening: %v", svr.name, server.ListenAndServe())
		}
	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop

	logger.Infof("%s server shutdown initiated...", svr.name)
	ctx, cancel := context.WithTimeout(context.Background(), svr.timeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("%s server shutdown error: %v", svr.name, err)
	}
	logger.Infof("%s server shutdown complete", svr.name)
}
