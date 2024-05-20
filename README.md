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
PRIVATE_KEY=YOUR_PRIVATE_KEY
RPC_URL=https://mainnet.aurora.dev
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
curl -X POST -H "Content-Type: application/json" -d '{"address":"0xRecipientAddress", "amount":"0.0001"}' http://localhost:8080/withdraw
```

address (string, required): 接收转账的以太坊地址。
amount (float64, required): 转账的金额，以 ETH 为单位。

## 响应案例

```json
{
    "hash": "0xTransactionHash"
}
```

```json
{
    "error": "xxxxxxxxxxxx"
}
```


## 协议和签名方法

```plaintext
LegacyTx协议
NewLondonSigner签名方法
```

## 注意事项

```plaintext
确保 `.env` 文件中的私钥和 RPC URL 正确无误。
如果使用的是本地开发链，请确保该链已经启动，并且 RPC URL 指向正确的地址。
```