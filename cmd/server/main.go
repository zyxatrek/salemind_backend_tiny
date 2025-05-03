package main

import (
	"log"

	"salemind_backend_tiny/pkg/api"
	"salemind_backend_tiny/pkg/config"
	"salemind_backend_tiny/pkg/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 创建服务
	genSvc := services.NewGenerationService(cfg)
	handler := api.NewHandler(genSvc)

	// 创建Gin引擎
	r := gin.Default()

	// 注册路由
	r.POST("/api/image_task/create", handler.CreateImageTask)
	r.POST("/api/image_task/status", handler.GetImageTaskStatus)
	r.POST("/api/video_task/create", handler.CreateVideoTask)
	r.POST("/api/video_task/status", handler.GetVideoTaskStatus)

	// 启动服务器
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
