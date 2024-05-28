# Aurora Faucet

Aurora Faucet 是一个使用 Go 语言编写的水龙头服务，允许用户请求一定数量的 Aurora (以太币) 通过 HTTP API。

## 目录结构

```plaintext
aurora-faucet/
├── .env
├── .gitignore
├── Makefile
├── README.md
├── go.mod
├── go.sum
├── build/
├── cmd/
│   └── aurora-faucet/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handlers/
│   │   └── withdraw.go
│   ├── services/
│   │   └── eth.go
│   ├── utils/
│   │   └── response.go

```

## 环境要求

```plaintext
Go 1.16 或更高版本
```

## 安装

克隆仓库并进入目录：

```bash
git clone git@github.com:lonySp/go-aurora-faucet.git
cd go-aurora-faucet
```

创建 .env 文件并添加以下内容：

```plaintext
PRIVATE_KEY=YOUR_PRIVATE_KEY                   // 转账地址私钥
PORT=8080                                      // 启动端口
AURORA_TESTNET_RPC_URL=https://testnet.aurora.dev   // 测试网络地址（目前通过前端传参，但先放着）
AURORA_TESTNET_EXPLORER_URL=https://testnet.aurorascan.dev  // 测试网络浏览器地址

```

安装依赖：

```bash
make install
```

## 使用

## 构建项目

```bash
make build
```

## 运行项目

```bash
make run
```

## 发送请求

使用 `curl` 或 Postman 发送请求：

```bash
curl -X POST -H "Content-Type: application/json" -d '{"address":"0xRecipientAddress", "network":"https://testnet.aurora.dev"}' http://localhost:8080/request
```

```bash
{
  "address": "xxxxxxx",
  "network": "https://testnet.aurora.dev",
  "amount": 0.0001
}
```



address (string, required): 接收ETH地址。\
network (string, required): 接收ETH网络地址 \
amount (float64, required): 提取金额

## 响应案例

```json
{
  "explorer_url": "https://testnet.aurorascan.dev/tx/xxxxxx",
  "success": true,
  "tx_id": "xxxxx"
}
```

```json
{
  "message": "余额不足",
  "success": false
}
```


## 协议和签名方法

```plaintext
LegacyTx协议
NewLondonSigner签名方法
```

## 注意事项

```plaintext
确保 `.env` 文件中的信息无误
```