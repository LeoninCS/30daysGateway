package core

import (
	"gateway/internal/config"
	"strings"
)

type Router struct {
	routes []*config.Route
}

func NewRouter(routes []*config.Route) *Router {
	return &Router{
		routes: routes,
	}
}

func (r *Router) FindRoute(path string) *config.Route {
	// 遍历查找匹配的路由
	for _, route := range r.routes {
		// 检查路径是否匹配
		if strings.HasPrefix(path, route.Path) {
			return route
		}
	}
	return nil
}
