package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"main/token"
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
	balance, err := client.BalanceAt(context.Background(), address, nil) // nil big.NewInt(5532992)
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

	// 代币token 余额
	// Golem (GNT) Address

	//solc --abi -o . erc20.sol
	// go get -u github.com/ethereum/go-ethereum/cmd/abigen
	//abigen --abi=ERC20.abi --pkg=token --out=erc20.go

	// 0xa74476443119A942dE498590Fe1f2454d7D4aC0d 是合约地址
	GNT_CONTRACT_ADDR := "0xa74476443119A942dE498590Fe1f2454d7D4aC0d"
	GNX_CONTRACT_ADDR := "0x6ec8a24cabdc339a06a172f8223ea557055adaa5"
	_ = GNT_CONTRACT_ADDR
	contractAddr := GNX_CONTRACT_ADDR
	tokenContractAddr := common.HexToAddress(contractAddr)
	newToken, err := token.NewToken(tokenContractAddr, client)
	if err != nil {
		fmt.Println(err)
		return
	}

	address = common.HexToAddress("0x0536806df512d6cdde913cf95c9886f65b1d3462")
	tokenBalance, err := newToken.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		fmt.Println(err)
		return
	}

	name, err := newToken.Name(&bind.CallOpts{})
	if err != nil {
		fmt.Println(err)
		return
	}

	symbol, err := newToken.Symbol(&bind.CallOpts{})
	decimal, err := newToken.Decimals(&bind.CallOpts{})

	fbal := new(big.Float)
	fbal.SetString(tokenBalance.String())
	bigFBalance := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimal))))

	fmt.Println("name: ", name)
	fmt.Println("symbol: ", symbol)
	fmt.Println("decimal: ", decimal)
	fmt.Println("tokenBalance: ", tokenBalance)
	fmt.Println("bigFBalance: ", bigFBalance)
}
