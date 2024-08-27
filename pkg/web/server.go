package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-pkgz/lgr"
)

func New(log lgr.L) *Server {
	return &Server{
		log: log,
	}
}

type Server struct {
	log lgr.L
	srv *http.Server
}

func (s *Server) Run(cancel context.CancelCauseFunc, port int, h http.Handler) {
	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: h,
	}

	go func() {
		err := s.srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			s.log.Logf("[ERROR] server failed: %v", err)
			cancel(err)
		} else {
			s.log.Logf("[INFO] server stopped")
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.srv == nil {
		return nil
	}

	err := s.srv.Shutdown(ctx)
	if err != nil {
		s.log.Logf("[ERROR] server shutdown failed: %v", err)
	} else {
		s.log.Logf("[INFO] server shutdown")
	}

	return err
}
