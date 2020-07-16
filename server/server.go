package server

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	Addr    string
	Handler http.Handler
}

func New(address string, handler http.Handler) *Server {
	return &Server{
		Addr:    address,
		Handler: handler,
	}
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	var g errgroup.Group
	s1 := &http.Server{
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       1 * time.Minute,
		WriteTimeout:      2 * time.Minute,
		IdleTimeout:       120 * time.Second,
		Addr:              s.Addr,
		Handler:           s.Handler,
	}

	g.Go(func() error {
		select {
		case <-ctx.Done():
			return s1.Shutdown(ctx)
		}
	})

	g.Go(func() error {
		return s1.ListenAndServe()
	})

	return g.Wait()
}
