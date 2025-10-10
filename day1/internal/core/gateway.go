package core

import (
	"gateway/internal/config"
	"net/http"
)

type Gateway struct {
	config *config.Config
	router *Router
	proxy  *Proxy
}

func NewGateway(config *config.Config) *Gateway {
	return &Gateway{
		config: config,
		router: NewRouter(config.Routes),
		proxy:  NewProxy(),
	}
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 查找匹配的路由
	route := g.router.FindRoute(r.URL.Path)
	if route == nil {
		http.Error(w, "路由未找到", http.StatusNotFound)
		return
	}

	// 转发请求
	g.proxy.ServeHTTP(w, r, route)
}
