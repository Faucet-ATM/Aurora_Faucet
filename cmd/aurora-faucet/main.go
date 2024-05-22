package main

import (
	"go-aurora-faucet/internal/config"
	"go-aurora-faucet/internal/handlers"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 创建共享的map和mutex
	lastWithdrawals := make(map[string]time.Time)
	mutex := &sync.Mutex{}

	// 创建提取处理器
	withdrawHandler := handlers.NewWithdrawHandler(cfg, lastWithdrawals, mutex)

	// 初始化 Gin 路由
	router := gin.Default()
	router.POST("/request", withdrawHandler.Handle)

	// 启动服务器
	log.Fatal(router.Run(":8080"))
}
