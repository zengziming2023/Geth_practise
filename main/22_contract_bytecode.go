package main

import (
	"context"
	"encoding/hex"
	"flag"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/glog"
	"main/main/constant"
)

func main() {
	flag.Parse()
	flag.Set("logtostderr", "true")
	flag.Set("v", "4") // 设置日志级别
	defer glog.Flush()

	client, err := ethclient.Dial(constant.SepoliaTest)
	if err != nil {
		glog.Fatal("dial err:", err)
	}

	contractAddress := common.HexToAddress(constant.StoreContractAddr)

	codeBytes, err := client.CodeAt(context.Background(), contractAddress, nil)
	if err != nil {
		glog.Fatal("codeAt err:", err)
	}

	glog.Info("codeAt:", hex.EncodeToString(codeBytes))

}
