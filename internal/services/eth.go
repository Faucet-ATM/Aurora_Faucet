package services

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ETHService struct {
	client          *ethclient.Client
	privateKey      *ecdsa.PrivateKey
	withdrawLimit   time.Duration
	lastWithdrawals map[string]time.Time
	mutex           *sync.Mutex
}

// NewETHService 创建ETHService实例
func NewETHService(privateKeyHex, networkURL string, withdrawLimit time.Duration, lastWithdrawals map[string]time.Time, mutex *sync.Mutex) (*ETHService, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("无法加载私钥: %v", err)
	}

	client, err := ethclient.Dial(networkURL)
	if err != nil {
		return nil, fmt.Errorf("无法连接到以太坊客户端: %v", err)
	}

	return &ETHService{
		client:          client,
		privateKey:      privateKey,
		withdrawLimit:   withdrawLimit,
		lastWithdrawals: lastWithdrawals,
		mutex:           mutex,
	}, nil
}

// CanWithdraw 检查用户是否可以提取
func (s *ETHService) CanWithdraw(userAddress string) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	lastWithdrawal, exists := s.lastWithdrawals[userAddress]
	if exists && time.Since(lastWithdrawal) < s.withdrawLimit {
		return false, fmt.Errorf("24小时内已经领取过")
	}
	return true, nil
}

// RecordWithdrawal 记录用户的提取时间
func (s *ETHService) RecordWithdrawal(userAddress string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.lastWithdrawals[userAddress] = time.Now()
}

// SendETH 发送ETH到指定地址
func (s *ETHService) SendETH(toAddress string, amount *big.Int) (string, error) {
	publicKey := s.privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKey)
	nonce, err := s.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", fmt.Errorf("无法获取 nonce: %v", err)
	}

	gasLimit := big.NewInt(21000)
	gasPrice, err := s.client.SuggestGasPrice(context.Background())
	if err != nil {

		return "", fmt.Errorf("无法获取 gasPrice: %v", err)
	}

	balance, err := s.client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		return "", fmt.Errorf("无法获取账户余额: %v", err)
	}

	totalCast := new(big.Int).Add(amount, new(big.Int).Mul(gasPrice, gasLimit))
	if balance.Cmp(totalCast) < 0 {
		return "", fmt.Errorf("余额不足")
	}

	// 创建交易
	to := common.HexToAddress(toAddress)
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    amount,
		Gas:      21000,
		GasPrice: gasPrice,
	})
	chainId, err := s.client.NetworkID(context.Background())
	if err != nil {
		return "", fmt.Errorf("获取链 ID 失败: %v", err)
	}

	// 签署交易
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainId), s.privateKey)
	if err != nil {
		return "", fmt.Errorf("签署交易失败: %v", err)
	}
	// 发送交易
	if err := s.client.SendTransaction(context.Background(), signedTx); err != nil {
		return "", fmt.Errorf("发送交易失败: %v", err)
	}
	return signedTx.Hash().Hex(), nil
}
