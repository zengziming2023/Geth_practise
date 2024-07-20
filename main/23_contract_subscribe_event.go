package main

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/glog"
	"main/main/constant"
	"main/main/util"
)

func main() {
	util.GlogInit()
	defer glog.Flush()

	client, err := ethclient.Dial(constant.SepoliaTestRPC)
	if err != nil {
		glog.Fatal("dial err:", err)
	}
	glog.Info("client:", client)
	contractAddr := common.HexToAddress(constant.StoreContractAddr)
	_ = contractAddr
	query := ethereum.FilterQuery{Addresses: []common.Address{contractAddr}} //

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		glog.Fatal("subscribeFilterLogs err:", err)
	}

	for {
		select {
		case err := <-sub.Err():
			glog.Fatal("err:", err)
		case vLog := <-logs:
			glog.Info("receive log:", vLog)
		}
	}

}
