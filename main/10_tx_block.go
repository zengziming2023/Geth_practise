package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

func main() {
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		return
	}
	// 查询最新的区域头
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return
	}
	fmt.Printf("header: %v\n", header.Number.String())
	num := big.NewInt(5671744)
	block, err := client.BlockByNumber(context.Background(), num)
	if err != nil {
		return
	}
	fmt.Printf("block: %v\n", block.Number().String())
	fmt.Printf("hash: %s\n", block.Hash().Hex())
	fmt.Printf("difficulty: %s\n", block.Difficulty().String())
	fmt.Printf("nonce: %d\n", block.Nonce())
	fmt.Printf("gasLimit: %d\n", block.GasLimit())
	fmt.Printf("gasUsed: %d\n", block.GasUsed())
	fmt.Printf("hash: %s\n", block.Hash().Hex())
	fmt.Printf("Transactions: %d\n", len(block.Transactions()))

	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		return
	}
	fmt.Printf("number of transactions: %d\n", count)
}
