package core

import (
	"gateway/internal/config"
	"gateway/internal/core/loadbalancer"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Proxy struct {
	loadbalancers map[string]loadbalancer.LoadBalancer
}

func NewProxy(backends []*config.Route) *Proxy {
	loadbalancers := make(map[string]loadbalancer.LoadBalancer)
	for _, backend := range backends {
		var lb loadbalancer.LoadBalancer
		switch backend.Strategy {
		case "round":
			lb = loadbalancer.NewRoundRobinLoadBalancer()
		case "random":
			lb = loadbalancer.NewRandomLoadBalancer()
		case "leastConn":
			lb = loadbalancer.NewLeastConnLoadBalancer()
		case "weight":
			lb = loadbalancer.NewWeightedLoadBalancer()
		default:
			lb = loadbalancer.NewRoundRobinLoadBalancer()
		}
		for _, target := range backend.Targets {
			backendURL, _ := url.Parse(target.URL)
			lb.AddBackend(&loadbalancer.Backend{
				ID:     backend.ID,
				URL:    backendURL,
				Alive:  true,
				Weight: target.Weight,
			})
		}
		loadbalancers[backend.ID] = lb
	}
	return &Proxy{
		loadbalancers: loadbalancers,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request, route *config.Route) {
	backend := p.loadbalancers[route.ID].GetNextTarget()
	if backend == nil {
		http.Error(w, "无可用后端服务", http.StatusServiceUnavailable)
		return
	}
	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(backend.URL)

	// 保存原始路径用于日志记录
	originalPath := r.URL.Path

	// 设置反向代理的Director
	proxy.Director = func(req *http.Request) {
		// 设置目标服务的协议和主机
		req.URL.Scheme = backend.URL.Scheme
		req.URL.Host = backend.URL.Host
		// 处理路径映射
		req.URL.Path = backend.URL.Path + strings.TrimPrefix(req.URL.Path, route.Path)
	}

	// 记录转发日志
	log.Printf("转发请求: %s %s -> %s://%s%s", r.Method, originalPath, backend.URL.Scheme, backend.URL.Host, r.URL.Path)

	// 转发请求
	proxy.ServeHTTP(w, r)
}
