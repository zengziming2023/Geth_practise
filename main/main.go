package main

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"time"
)

func main() {
	t := time.Now()
	//client, err := ethclient.Dial("https://cloudflare-eth.com") 	// 远端

	// npm install -g ganache-cli 	安装
	// ganache-cli 	运行
	// http://localhost:8584
	client, err := ethclient.Dial("http://localhost:8584") // 本地
	if err != nil {
		log.Fatal(err)
		return
	}
	cost := time.Since(t)
	log.Printf("cost: %v\n", cost)
	log.Println("Connected to Ethereum client")
	_ = client
}
