package router

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"maps"

	httpdrv "github.com/mirkobrombin/go-module-router/v1/http"
	"github.com/mirkobrombin/go-module-router/v1/logger"
	"github.com/mirkobrombin/go-module-router/v1/middleware"
	"github.com/mirkobrombin/go-module-router/v1/registry"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

func Build(reg *registry.Registry) []registry.Route {
	var out []registry.Route
	for _, rp := range reg.RouteProviders {
		out = append(out, rp()...)
	}
	return out
}

func New(reg *registry.Registry, services map[string]any, eng httpdrv.Engine, opt Options) httpdrv.Engine {
	if opt.Logger == nil {
		opt.Logger = logger.Nop
	}

	repoInstances := map[string]any{}
	maps.Copy(repoInstances, services)

	if !opt.SkipAutoWire {
		for name, fn := range reg.RepoInit {
			repoInstances[name] = fn(repoInstances)
		}
	}

	svcInstances := map[string]any{}
	maps.Copy(svcInstances, services)
	if !opt.SkipAutoWire {
		for name, fn := range reg.ServiceInit {
			if _, ok := svcInstances[name]; !ok {
				svcInstances[name] = fn(repoInstances)
			}
		}
	}

	if eng == nil {
		eng = httpdrv.NewFastHTTP()
	}

	mwInstances := map[string]any{}
	for name, fn := range reg.MiddlewareInit {
		mwInstances[name] = fn(svcInstances, opt.SessionDuration)
	}

	hInstances := map[string]any{}
	for name, fn := range reg.HandlerInit {
		hInstances[name] = fn(svcInstances)
	}

	for _, rt := range Build(reg) {
		parts := strings.Split(rt.HandlerName, ".")
		module := parts[0]

		if v := os.Getenv("GMR_MOD_OFF_" + strings.ToUpper(module)); v != "" {
			if disabled, err := strconv.ParseBool(v); err == nil && disabled {
				opt.Logger.Info("router: skipping entire module", "module", module)
				continue
			}
		}

		handlerKey := strings.ToUpper(strings.ReplaceAll(rt.HandlerName, ".", "_"))
		if v := os.Getenv("GMR_HAND_OFF_" + handlerKey); v != "" {
			if disabled, err := strconv.ParseBool(v); err == nil && disabled {
				opt.Logger.Info("router: skipping handler", "handler", rt.HandlerName)
				continue
			}
		}

		h, err := resolveHandler(hInstances, rt.HandlerName)
		if err != nil {
			opt.Logger.Error("router: handler reflection failed",
				zap.String("handler", rt.HandlerName),
				zap.Error(err))
			if opt.OnError != nil {
				opt.OnError(err)
			}
			continue
		}

		fh := h
		for i := len(rt.Middleware) - 1; i >= 0; i-- {
			raw := mwInstances[rt.Middleware[i]]
			mw, ok := raw.(middleware.Component)
			if !ok {
				opt.Logger.Error("router: middleware does not implement Component",
					zap.String("name", rt.Middleware[i]))
				continue
			}
			fh = mw.Apply(fh, middleware.RouteInfo{Permissions: rt.Permissions})
		}

		eng.Handle(rt.Method, rt.Path, fh)
		opt.Logger.Debug("router: registered",
			zap.String("method", rt.Method),
			zap.String("path", rt.Path),
			zap.String("handler", rt.HandlerName))
	}

	return eng
}

func resolveHandler(hmap map[string]any, fq string) (fasthttp.RequestHandler, error) {
	parts := strings.Split(fq, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid handler name %q", fq)
	}
	inst, ok := hmap[parts[0]]
	if !ok {
		return nil, fmt.Errorf("handler %s not found", parts[0])
	}
	m := reflect.ValueOf(inst).MethodByName(parts[1])
	if !m.IsValid() {
		return nil, fmt.Errorf("method %s not found on %T", parts[1], inst)
	}
	t := m.Type()
	if t.NumIn() != 1 ||
		t.In(0) != reflect.TypeOf(&fasthttp.RequestCtx{}) ||
		t.NumOut() != 0 {
		return nil, fmt.Errorf("method %s has incompatible signature", fq)
	}
	return func(ctx *fasthttp.RequestCtx) {
		m.Call([]reflect.Value{reflect.ValueOf(ctx)})
	}, nil
}
