package main

import (
	"gateway/internal/config"
	"gateway/internal/core"
	"log"
)

func main() {
	// 加载配置
	log.Println("开始加载配置")
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatal("加载配置失败:", err)
	}
	log.Println("加载配置成功")

	// 创建网关
	log.Println("开始创建网关")
	gw := core.NewGateway(cfg)
	log.Println("网关创建成功")

	// 创建服务器
	log.Println("开始创建服务器")
	// 修复：确保端口地址格式正确
	addr := ":" + cfg.Server.Port
	srv, err := core.NewServer(addr, gw)
	if err != nil {
		log.Fatal("创建服务器失败:", err)
	}
	log.Println("服务器创建成功")

	// 启动服务器
	log.Printf("服务器正在监听端口 %s", cfg.Server.Port)
	if err := srv.Start(); err != nil {
		log.Fatal("服务器启动失败:", err)
	}
	log.Println("服务器已启动")
}
