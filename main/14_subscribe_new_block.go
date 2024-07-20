package main

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/glog"
	"main/main/constant"
	"main/main/util"
	"time"
)

func main() {
	util.GlogInit()
	defer glog.Flush()

	// websocket RPC https://sepolia.infura.io/v3/e732e9f17ce2413c884fa5b4a6960ee3
	client, err := ethclient.Dial(constant.SepoliaTestRPC)
	if err != nil {
		glog.Fatal("dial: ", err)
	}

	//创建一个新的通道，用于接收最新的区块头
	headers := make(chan *types.Header)

	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		glog.Fatal("subscribe: ", err)
	}

	for {
		select {
		case err := <-sub.Err():
			glog.Fatal("sub error: ", err)
		case header := <-headers:
			//glog.Info(header.Hash().Hex())
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				glog.Fatal("block by hash: ", err)
			}

			glog.Info("==== receive block start ====")
			glog.Info(block.Hash().Hex())
			glog.Info(block.Number().Uint64())
			glog.Info(time.Unix(int64(block.Time()), 0).Local().Format("2006-01-02 15:04:05"))
			glog.Info(block.Nonce())
			glog.Info(len(block.Transactions()))
			glog.Info("==== receive block end ====\n")
		}
	}
}
