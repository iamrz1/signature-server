package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"signature-server/util"
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
	util.Infof("starting %s server...", svr.name)

	server := http.Server{
		Addr:              fmt.Sprintf(":%d", svr.port),
		Handler:           svr.handler,
		ReadTimeout:       svr.timeout,
		ReadHeaderTimeout: svr.timeout,
		WriteTimeout:      svr.timeout,
		IdleTimeout:       svr.timeout,
	}

	go func() {
		util.Infof("%s server listening on port %d", svr.name, svr.port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			util.Errorf("%s server stopped listening: %v", svr.name, server.ListenAndServe())
		}
	}()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop

	util.Infof("%s server shutdown initiated...", svr.name)
	ctx, cancel := context.WithTimeout(context.Background(), svr.timeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		util.Errorf("%s server shutdown error: %v", svr.name, err)
	}
	util.Infof("%s server shutdown complete", svr.name)
}
