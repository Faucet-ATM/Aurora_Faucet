package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	PrivateKey string
	// WithdrawAmount           *big.Int
	Port string
	// WithdrawLimit            time.Duration
	AuroraTestnetRPCURL      string
	AuroraTestnetExplorerURL string
}

func LoadConfig() *Config {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// 解析提取金额
	//withdrawAmountETH, err := strconv.ParseFloat(os.Getenv("WITHDRAW_AMOUNT"), 64)
	//if err != nil {
	//	log.Fatalf("Invalid WITHDRAW_AMOUNT: %v", err)
	//}
	//withdrawAmountWei := new(big.Int).Mul(big.NewInt(int64(withdrawAmountETH*1e18)), big.NewInt(1))

	// 解析提取时间限制
	//withdrawLimit, err := strconv.Atoi(os.Getenv("WITHDRAW_LIMIT"))
	//if err != nil {
	//	log.Fatalf("Invalid WITHDRAW_LIMIT: %v", err)
	//}
	return &Config{
		PrivateKey: os.Getenv("PRIVATE_KEY"),
		Port:       os.Getenv("PORT"),
		// WithdrawLimit: time.Duration(withdrawLimit) * time.Hour,
		// WithdrawAmount:           withdrawAmountWei,
		AuroraTestnetRPCURL:      os.Getenv("AURORA_TESTNET_RPC_URL"),
		AuroraTestnetExplorerURL: os.Getenv("AURORA_TESTNET_EXPLORER_URL"),
	}
}
