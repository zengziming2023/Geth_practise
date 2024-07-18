package main

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func main() {
	client, err := ethclient.Dial("https://cloudflare-eth.com")
	if err != nil {
		log.Println(err)
		return
	}

	blockNum := big.NewInt(5671744)
	block, err := client.BlockByNumber(context.Background(), blockNum)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(len(block.Transactions()))
	for index, tx := range block.Transactions()[0:5] {
		log.Println(index, tx)
		log.Println(tx.Hash().Hex())
		log.Println(tx.Value().String())
		log.Println(tx.Gas())
		log.Println(tx.GasPrice().Uint64())
		log.Println(tx.Nonce())
		log.Println(tx.Data())
		log.Println(tx.To().Hex())

		networkID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Println(err)
		}

		// get sender info
		sender, err := types.Sender(types.NewEIP155Signer(networkID), tx)
		if err != nil {
			log.Println(err)
		}
		log.Println(sender.Hex())

		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Println(err)
		}
		log.Println(receipt.Status)
	}

	log.Println("======")
	blockHash := common.HexToHash("0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9")
	count, err := client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		log.Println(err)
	}

	log.Println(count)
	for idx := range count {
		transaction, err := client.TransactionInBlock(context.Background(), blockHash, idx)
		if err != nil {
			log.Println(err)
		}
		log.Println(transaction.Hash().Hex())
		if idx > 1 {
			break
		}
	}

	txHash := common.HexToHash("0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2")
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Println(err)
	}
	log.Println(isPending)
	log.Println(tx.Hash().Hex())
}
