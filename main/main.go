package main

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"time"
)

func main() {
	t := time.Now()
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Fatal(err)
		return
	}
	cost := time.Since(t)
	log.Printf("cost: %v\n", cost)
	log.Println("Connected to Ethereum client")
	_ = client
}
