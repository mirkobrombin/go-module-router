package http

import (
	"context"
	"net/http"
)

type StdHTTP struct {
	mux     *http.ServeMux
	mwChain []Middleware
	srv     *http.Server
}

func NewStdHTTP() *StdHTTP {
	return &StdHTTP{mux: http.NewServeMux()}
}

func (s *StdHTTP) Handle(m, p string, h Handler) {
	final := h.(http.Handler)
	for i := len(s.mwChain) - 1; i >= 0; i-- {
		final = s.mwChain[i](final).(http.Handler)
	}

	switch m {
	case http.MethodGet:
		s.mux.Handle(p, final)
	default:
		s.mux.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				http.NotFound(w, r)
				return
			}
			final.ServeHTTP(w, r)
		})
	}
}

func (s *StdHTTP) Use(mw Middleware) { s.mwChain = append(s.mwChain, mw) }

func (s *StdHTTP) Group(prefix string) Engine {
	sub := NewStdHTTP()
	sub.mwChain = s.mwChain
	s.mux.Handle(prefix+"/", http.StripPrefix(prefix, sub.mux))
	return sub
}

func (s *StdHTTP) Serve(addr string) error {
	s.srv = &http.Server{Addr: addr, Handler: s.mux}
	return s.srv.ListenAndServe()
}

func (s *StdHTTP) Shutdown(ctx context.Context) error {
	if s.srv == nil {
		return nil
	}
	return s.srv.Shutdown(ctx)
}

func (s *StdHTTP) Unwrap() any { return s.mux }
