package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
	"main/token"
	"math"
	"math/big"
	"os"
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

	// 生成新钱包
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println(err)
		return
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println("privateKeyHexStr: ", hexutil.Encode(privateKeyBytes)[2:])

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("error casting public key to ECDSA")
		return
	}

	// 剥离了 0x 和前 2 个字符 04，它始终是 EC 前缀，不是必需的
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("publicKeyHexStr: ", hexutil.Encode(publicKeyBytes)[4:])
	// public key 转 public address
	pubAddress := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("pubAddress: ", pubAddress)

	// 公共地址其实就是公钥的 Keccak-256 哈希，然后我们取最后 40 个字符（20 个字节）并用“0x”作为前缀。 以下是使用 golang.org/x/crypto/sha3 的 Keccak256 函数手动完成的方法。
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println("hash", hexutil.Encode(hash.Sum(nil)[12:]))

	//createKs()
	importKs()
}

func createKs() {
	ks := keystore.NewKeyStore("./keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	password := "password"
	account, err := ks.NewAccount(password)
	if err != nil {
		return
	}
	fmt.Println("account: ", account.Address.Hex())
}

func importKs() {
	file := "./keystore/UTC--2024-07-17T10-17-58.435896000Z--fb770dad18a7b4044de3b7db33b3e4882a10400e"
	ks := keystore.NewKeyStore("./keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		//return
	}
	password := "password"
	account, err := ks.Import(jsonBytes, password, password)
	if err != nil {
		fmt.Println(err)
		//return
	}
	fmt.Println("account: ", account.Address.Hex()) // should be 0xfB770dad18A7b4044de3b7db33b3e4882A10400e

}
