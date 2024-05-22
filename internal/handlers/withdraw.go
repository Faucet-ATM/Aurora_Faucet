package handlers

import (
	"go-aurora-faucet/internal/config"
	"go-aurora-faucet/internal/services"
	"go-aurora-faucet/internal/utils"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// WithdrawRequest 结构体用于接收提取请求的参数
type WithdrawRequest struct {
	Network string `json:"network" binding:"required"` // 网络节点 URL
	Address string `json:"address" binding:"required"` // 目标钱包地址
}

// WithdrawHandler 处理提取请求的处理器
type WithdrawHandler struct {
	config          *config.Config
	lastWithdrawals map[string]time.Time
	mutex           *sync.Mutex
}

// NewWithdrawHandler 创建一个新的 WithdrawHandler 实例
func NewWithdrawHandler(config *config.Config, lastWithdrawals map[string]time.Time, mutex *sync.Mutex) *WithdrawHandler {
	return &WithdrawHandler{
		config:          config,
		lastWithdrawals: lastWithdrawals,
		mutex:           mutex,
	}
}

// Handle 处理提取请求
func (h *WithdrawHandler) Handle(c *gin.Context) {
	var req WithdrawRequest
	// 绑定请求参数到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 动态创建 ETH 服务
	ethService, err := services.NewETHService(h.config.PrivateKey, req.Network, h.config.WithdrawLimit, h.lastWithdrawals, h.mutex)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Failed to create ETH service: "+err.Error())
		return
	}

	// 检查是否可以领取
	canWithdraw, err := ethService.CanWithdraw(req.Address)
	if err != nil {
		utils.RespondError(c, http.StatusTooManyRequests, err.Error())
		return
	}
	if !canWithdraw {
		utils.RespondError(c, http.StatusTooManyRequests, "24小时内只能领取一次")
		return
	}

	// 发送交易
	txHash, err := ethService.SendETH(req.Address, h.config.WithdrawAmount)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 记录用户提取时间
	ethService.RecordWithdrawal(req.Address)

	// 返回交易哈希和区块浏览器 URL
	explorerURL := h.config.AuroraTestnetExplorerURL + "/tx/" + txHash
	utils.RespondSuccess(c, gin.H{
		"tx_id":        txHash,
		"explorer_url": explorerURL,
	})
}
