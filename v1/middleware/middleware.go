package middleware

import "github.com/valyala/fasthttp"

type RouteInfo struct {
	Permissions []string
}

type Component interface {
	Apply(next fasthttp.RequestHandler, info RouteInfo) fasthttp.RequestHandler
}
