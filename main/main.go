package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math"
	"math/big"
	"time"
)

func main() {
	t := time.Now()
	client, err := ethclient.Dial("https://cloudflare-eth.com") // 远端

	// npm install -g ganache-cli 	安装
	// ganache-cli 	运行
	// http://localhost:8584
	//client, err := ethclient.Dial("http://localhost:8584") // 本地
	if err != nil {
		fmt.Println(err)
		return
	}
	cost := time.Since(t)
	fmt.Printf("cost: %v\n", cost)
	fmt.Println("Connected to Ethereum client")
	//_ = client

	address := common.HexToAddress("0x71c7656ec7ab88b098defb751b7401b5f6d8976f")
	fmt.Println("address: ", address)
	fmt.Println("address: ", address.Hex())
	fmt.Println("address: ", address.Bytes())

	// nil 表示最新区块
	balance, err := client.BalanceAt(context.Background(), address, big.NewInt(5532992)) // nil
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("balance: ", balance)

	blockNum := big.NewInt(5532993)
	balance, err = client.BalanceAt(context.Background(), address, blockNum)
	fmt.Println("blockNum: ", blockNum, ", balance: ", balance)

	if balance == nil {
		fmt.Println("balance is nil")
		return
	}

	fBalance := new(big.Float)
	fBalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fBalance, big.NewFloat(math.Pow10(18)))
	fmt.Println("ethValue: ", ethValue)

	// 有时您想知道待处理的账户余额是多少，例如，在提交或等待交易确认后
	pendingBalanceAt, err := client.PendingBalanceAt(context.Background(), address)
	if err != nil {
		return
	}
	fmt.Println("pendingBalanceAt: ", pendingBalanceAt)
}
