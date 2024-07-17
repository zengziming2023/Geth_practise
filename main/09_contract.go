package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"regexp"
)

func main() {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	fmt.Println(re.MatchString("0x323b5d4c32345ced77393b3530b1eed0f346429d"))
	fmt.Println(re.MatchString("0xZYXb5d4c32345ced77393b3530b1eed0f346429d"))

	// 若在该地址存储了字节码，该地址是智能合约
	client, err := ethclient.Dial("https://cloudflare-eth.com") // 远端
	if err != nil {
		fmt.Println(err)
	}

	address := common.HexToAddress("0xe41d2489571d322189246dafa5ebde1f4699f498")
	codeBytes, err := client.CodeAt(context.Background(), address, nil)
	if err != nil {
		return
	}
	fmt.Println(string(codeBytes))
	isContract := len(codeBytes) > 0
	fmt.Println("isContract: ", isContract)
}
