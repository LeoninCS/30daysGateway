package core

import (
	"gateway/internal/config"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Proxy struct {
}

func NewProxy() *Proxy {
	return &Proxy{}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request, route *config.Route) {
	// 解析目标服务URL
	targetURL, err := url.Parse(route.Target)
	if err != nil {
		http.Error(w, "无效的目标服务地址", http.StatusInternalServerError)
		return
	}

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// 保存原始路径用于日志记录
	originalPath := r.URL.Path

	// 设置反向代理的Director
	proxy.Director = func(req *http.Request) {
		// 设置目标服务的协议和主机
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		// 处理路径映射
		req.URL.Path = targetURL.Path + strings.TrimPrefix(req.URL.Path, route.Path)
	}

	// 记录转发日志
	log.Printf("转发请求: %s %s -> %s://%s%s", r.Method, originalPath, targetURL.Scheme, targetURL.Host, r.URL.Path)

	// 转发请求
	proxy.ServeHTTP(w, r)
}
