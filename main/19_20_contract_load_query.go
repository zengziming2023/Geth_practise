package main

import (
	"flag"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/glog"
	"main/main/constant"
	"main/store"
)

func main() {
	flag.Parse()
	flag.Set("logtostderr", "true")
	flag.Set("v", "4") // 设置日志级别
	defer glog.Flush()

	client, err := ethclient.Dial(constant.SepoliaTest)
	if err != nil {
		glog.Fatal("dial eth err:", err)
	}

	storeContractAddr := common.HexToAddress(constant.StoreContractAddr)

	storeInstance, err := store.NewStore(storeContractAddr, client)
	if err != nil {
		glog.Fatal("new store err:", err)
	}
	glog.Info("new store success.")
	_ = storeInstance

	version, err := storeInstance.Version(nil)
	if err != nil {
		glog.Fatal("version err:", err)
	}
	glog.Info("version: ", version)

}
