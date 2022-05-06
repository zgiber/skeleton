package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/julienschmidt/httprouter"
)

const (
	shutdownTimeout = 20 * time.Second
)

type transport struct {
	mux        *httprouter.Router
	httpServer *http.Server
	service    *service
}

func NewTransport(srv *service, addr string) *transport {
	mux := httprouter.New()
	return &transport{
		mux: mux,
		httpServer: &http.Server{
			Handler: mux,
			Addr:    addr,
		},
		service: srv,
	}
}

func (t *transport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.httpServer.Handler.ServeHTTP(w, r)
}

func (t *transport) Handle(method, pattern string, fn http.HandlerFunc) {
	t.mux.HandlerFunc(method, pattern, fn)
}

func (t *transport) Run() error {
	c := make(chan os.Signal, 1)
	err := make(chan error, 1)
	go func() {
		if transportErr := t.httpServer.ListenAndServe(); err != nil {
			err <- transportErr
		}
	}()

	ctx, cf := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cf()
	defer t.httpServer.Shutdown(ctx)

	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		log.Println("Interrupt received, shutting down")
		return nil
	case transportErr := <-err:
		return transportErr
	}
}
