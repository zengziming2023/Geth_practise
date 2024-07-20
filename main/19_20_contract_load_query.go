package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/glog"
	"main/main/constant"
	"main/main/util"
	"main/store"
)

func main() {
	util.GlogInit()
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
