package registry

import "sync"

var (
	global     = New()
	onceGlobal sync.Once
)

func Global() *Registry { return global }

type Registry struct {
	ServiceInit    map[string]ServiceInit
	RepoInit       map[string]RepoInit
	HandlerInit    map[string]HandlerInit
	MiddlewareInit map[string]MiddlewareInit
	ModelProviders []ModelProvider
	RouteProviders []RouteProvider
	mu             sync.RWMutex
}

func New() *Registry {
	return &Registry{
		ServiceInit:    map[string]ServiceInit{},
		RepoInit:       map[string]RepoInit{},
		HandlerInit:    map[string]HandlerInit{},
		MiddlewareInit: map[string]MiddlewareInit{},
	}
}

func RegisterService(name string, fn ServiceInit) {
	global.mu.Lock()
	defer global.mu.Unlock()
	global.ServiceInit[name] = fn
}
func RegisterRepository(name string, fn RepoInit) {
	global.mu.Lock()
	defer global.mu.Unlock()
	global.RepoInit[name] = fn
}
func RegisterHandler(name string, fn HandlerInit) {
	global.mu.Lock()
	defer global.mu.Unlock()
	global.HandlerInit[name] = fn
}
func RegisterMiddleware(name string, fn MiddlewareInit) {
	global.mu.Lock()
	defer global.mu.Unlock()
	global.MiddlewareInit[name] = fn
}
func RegisterModels(fn ModelProvider) {
	global.mu.Lock()
	defer global.mu.Unlock()
	global.ModelProviders = append(global.ModelProviders, fn)
}
func RegisterRoutes(fn RouteProvider) {
	global.mu.Lock()
	defer global.mu.Unlock()
	global.RouteProviders = append(global.RouteProviders, fn)
}
