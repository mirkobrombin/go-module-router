package http

import (
	"context"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type fastGroup struct{ g *router.Group }

func (fg *fastGroup) Handle(m, p string, h Handler) {
	fg.g.Handle(m, p, h.(fasthttp.RequestHandler))
}
func (fg *fastGroup) Use(Middleware)                 {}
func (fg *fastGroup) Group(prefix string) Engine     { return &fastGroup{g: fg.g.Group(prefix)} }
func (fg *fastGroup) Serve(string) error             { return nil }
func (fg *fastGroup) Shutdown(context.Context) error { return nil }
func (fg *fastGroup) Unwrap() any                    { return fg.g }

type FastHTTP struct {
	r       *router.Router
	mwChain []Middleware
	srv     *fasthttp.Server
}

func NewFastHTTP() *FastHTTP {
	return &FastHTTP{r: router.New()}
}

func (f *FastHTTP) Handle(m, p string, h Handler) {
	final := h.(fasthttp.RequestHandler)
	for i := len(f.mwChain) - 1; i >= 0; i-- {
		final = f.mwChain[i](final).(fasthttp.RequestHandler)
	}
	f.r.Handle(m, p, final)
}

func (f *FastHTTP) Use(mw Middleware) { f.mwChain = append(f.mwChain, mw) }

func (f *FastHTTP) Group(prefix string) Engine {
	return &fastGroup{g: f.r.Group(prefix)}
}

func (f *FastHTTP) Serve(addr string) error {
	f.srv = &fasthttp.Server{Handler: f.r.Handler}
	return f.srv.ListenAndServe(addr)
}

func (f *FastHTTP) Shutdown(ctx context.Context) error {
	if f.srv == nil {
		return nil
	}
	return f.srv.ShutdownWithContext(ctx)
}

func (f *FastHTTP) Unwrap() any { return f.r }
