package main

import (
	bzzclient "github.com/ethersphere/swarm/api/client"
	"github.com/golang/glog"
	"main/main/util"
)

func main() {
	// swarm 太久没有更新了。。先降级版本来适配
	// go get github.com/ethereum/go-ethereum@v1.9.2
	// go get github.com/ethersphere/go-sw3@v0.2.1
	// go mod tidy
	// 跑完以后，再升级回去 go get github.com/ethereum/go-ethereum@latest
	// go get github.com/ethersphere/go-sw3@latest
	// go mod tidy

	util.GlogInit()
	defer glog.Flush()

	// 启动swarm 服务
	//export BZZKEY=44b7B3293C829cF2b440a63Bc14DE902fEa53733
	//swarm --bzzaccount $BZZKEY

	client := bzzclient.NewClient("http://127.0.0.1:8500")

	file, err := bzzclient.Open("./swarm_test_upload.txt")
	if err != nil {
		glog.Fatal("open swarm test upload file err:", err)
	}

	manifestHash, err := client.Upload(file, "", false, false, false)
	if err != nil {
		glog.Fatal("upload swarm test upload file err:", err)
	}
	glog.Info("manifest hash: ", manifestHash) // 437553621d2c6839e8c966037d9677f8b13de411729f56fad5cbd6b2286c8e60
}
