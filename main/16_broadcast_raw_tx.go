package main

import (
	"encoding/hex"
	"flag"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/golang/glog"
	"main/main/constant"
)

func main() {
	flag.Parse()
	flag.Set("logtostderr", "true")
	flag.Set("v", "4") // 设置日志级别
	defer glog.Flush()

	glog.Info("Starting Eth Transfer...")
	defer glog.Info("Stopping Eth Transfer...")
	// "https://rinkeby.infura.io/v3/e732e9f17ce2413c884fa5b4a6960ee3"
	client, err := ethclient.Dial(constant.SepoliaTest)
	if err != nil {
		glog.Info("dial error: ", err)
	}
	_ = client

	rawTx := "f86e038414ef894d825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d872386f26fc10000808401546d72a0bb4bb8e04ef0b10ae23b1bab088671c8f05734784d69aad1aa7367f7acabbe86a04313a91b88a7ced20948ac0f7523197bca2d487af9a48045adbc84e37ee7efda"
	rawTxBytes, err := hex.DecodeString(rawTx)
	if err != nil {
		glog.Fatal("decode raw tx error: ", err)
	}

	tx := new(types.Transaction)
	err = rlp.DecodeBytes(rawTxBytes, tx)
	if err != nil {
		glog.Fatal("decode raw tx 2 error: ", err)
	}

	//err = client.SendTransaction(context.Background(), tx)
	//if err != nil {
	//	glog.Fatal("send raw tx error: ", err)
	//}
	glog.Info("raw tx sent: ", tx.Hash().Hex())
}
