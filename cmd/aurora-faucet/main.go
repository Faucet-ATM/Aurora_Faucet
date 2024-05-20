package main

import (
	"go-aurora-faucet/internal/config"
	"go-aurora-faucet/internal/handlers"
	"go-aurora-faucet/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()
	ethService, err := services.NewETHService(cfg)
	if err != nil {
		log.Fatalf("Failed to create ETH service: %v", err)
	}

	handler := handlers.NewHandler(ethService)

	router := gin.Default()
	router.POST("/withdraw", handler.Withdraw)
	log.Fatal(router.Run(":8080"))
}
