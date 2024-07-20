package main

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/glog"
	"main/main/util"
	"math/big"
)

func main() {
	util.GlogInit()
	defer glog.Flush()

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		glog.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)

	balance := new(big.Int)
	balance.SetString("10000000000000000000", 10) // 10 eth

	addr := auth.From
	glog.Info("addr: ", addr.Hex())
	// 创建一个创世账户并为其分配初始余额
	genesisAlloc := map[common.Address]types.Account{
		addr: {
			Balance: balance,
		},
	}

	blockGasLimit := uint64(4712388)
	client := backends.NewSimulatedBackend(genesisAlloc, blockGasLimit)

	fromAddr := auth.From
	glog.Info("fromAddr: ", fromAddr.Hex())
	// get nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		glog.Fatal("nonce error:", err)
	}
	glog.Info("nonce: ", nonce)

	value := big.NewInt(100_0000_0000_0000_0000) // 1 eth
	// gas limit
	gasLimit := uint64(21000)
	// gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		glog.Fatal("gas price error:", err)
	}
	glog.Info("gasPrice: ", gasPrice)

	toAddr := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")

	var data []byte

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddr,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		Data:     data,
		Value:    value,
	})
	glog.Info("tx: ", tx.Hash().Hex())

	chainId := big.NewInt(1337)

	signTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		glog.Fatal("sign tx error:", err)
	}
	glog.Info("signTx: ", signTx.Hash().Hex())

	err = client.SendTransaction(context.Background(), signTx)
	if err != nil {
		glog.Fatal("send tx error:", err)
	}

	glog.Info("send signTx: ", signTx.Hash().Hex())

	// 挖矿
	client.Commit()

	receipt, err := client.TransactionReceipt(context.Background(), signTx.Hash())
	if err != nil {
		glog.Fatal("get tx receipt error:", err)
	}
	if receipt == nil {
		glog.Fatal("receipt is nil. Forgot to commit?")
	}
	glog.Info("status: ", receipt.Status)
}
