package handlers

import (
	"go-aurora-faucet/internal/services"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
)

// WithdrawRequest struct holds the request payload for the withdraw endpoint
type WithdrawRequest struct {
	Address string  `json:"address" binding:"required"`
	Amount  float64 `json:"amount" binding:"required"`
}

// Handler struct holds the services used by the handlers
type Handler struct {
	ethService *services.ETHService
}

// NewHandler creates a new Handler with the given services
func NewHandler(ethService *services.ETHService) *Handler {
	return &Handler{ethService: ethService}
}

// Withdraw handles the withdraw request
func (h *Handler) Withdraw(c *gin.Context) {
	var req WithdrawRequest
	// Bind JSON payload to the request struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 转换金额为 wei
	amount := new(big.Int).Mul(big.NewInt(int64(req.Amount*1e18)), big.NewInt(1))

	// Call the ETH service to send the transaction
	txHash, err := h.ethService.SendETH(req.Address, amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the transaction hash in the response
	c.JSON(http.StatusOK, gin.H{"hash": txHash})
}
