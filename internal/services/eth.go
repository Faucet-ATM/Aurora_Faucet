package services

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"go-aurora-faucet/internal/config"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ETHService struct holds the client and private key for interacting with Ethereum
type ETHService struct {
	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey
}

// NewETHService creates a new ETHService with the given config
func NewETHService(cfg *config.Config) (*ETHService, error) {
	// Load private key from config
	privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("无法加载私钥: %v", err)
	}

	// Connect to Ethereum client
	client, err := ethclient.Dial(cfg.RPCUrl)
	if err != nil {
		return nil, fmt.Errorf("无法连接到以太坊客户端: %v", err)
	}

	// Return the ETHService struct
	return &ETHService{
		client:     client,
		privateKey: privateKey,
	}, nil
}

// SendETH sends ETH to the specified address
func (s *ETHService) SendETH(toAddress string, amount *big.Int) (string, error) {
	// Get the public key from the private key
	publicKey := s.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("无法转换公钥")
	}

	// Get the from address from the public key
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := s.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", fmt.Errorf("无法获取 nonce: %v", err)
	}

	// Get the balance of the from address
	balance, err := s.client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return "", fmt.Errorf("无法获取账户余额: %v", err)
	}

	gasLimit := uint64(21000) // in units
	gasPrice, err := s.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("无法获取 gas price: %v", err)
	}

	// Check if there is enough balance to cover the transaction
	totalCost := new(big.Int).Add(amount, new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit))))
	if balance.Cmp(totalCost) < 0 {
		return "", fmt.Errorf("余额不足，无法支付交易费用和转账金额: 需要 %s, 但只有 %s", totalCost.String(), balance.String())
	}

	// Create the transaction
	to := common.HexToAddress(toAddress)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: gasPrice,
	})

	chainID, err := s.client.NetworkID(context.Background())
	if err != nil {
		return "", fmt.Errorf("获取链 ID 失败: %v", err)
	}

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), s.privateKey)
	if err != nil {
		return "", fmt.Errorf("签署交易失败: %v", err)
	}

	// Send the transaction
	err = s.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", fmt.Errorf("发送交易失败: %v", err)
	}

	return signedTx.Hash().Hex(), nil
}
