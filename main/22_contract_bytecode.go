package main

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/glog"
	"main/main/constant"
	"main/main/util"
)

func main() {
	util.GlogInit()
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
