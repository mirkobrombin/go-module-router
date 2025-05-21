package middleware

import "github.com/mirkobrombin/go-module-router/v1/http"

type RouteInfo struct {
	Permissions []string
}

type Component interface {
	Apply(next http.Handler, info RouteInfo) http.Handler
}
