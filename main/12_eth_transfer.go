package main

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/glog"
	"main/main/constant"
	"main/main/util"
	"math"
	"math/big"
)

func main() {
	util.GlogInit()
	defer glog.Flush()

	glog.Info("Starting Eth Transfer...")
	defer glog.Info("Stopping Eth Transfer...")
	// "https://rinkeby.infura.io/v3/e732e9f17ce2413c884fa5b4a6960ee3"
	client, err := ethclient.Dial(constant.SepoliaTest)
	if err != nil {
		glog.Info("dial error: ", err)
	}

	// get privateKey
	privateKey, err := crypto.HexToECDSA(constant.PrivateKey)
	if err != nil {
		glog.Error("HexToECDSA error: ", err)
	}

	// get public key
	publicKey := privateKey.Public()
	// translate to ecdsa.PublicKey
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		glog.Error("error casting public key to ECDSA")
	}

	// get sender address
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// get balance
	balanceAt, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		glog.Error("balance error: ", err)
	}
	glog.Infof("BalanceAt: %v\n", balanceAt)
	fb := new(big.Float)
	fb.SetInt(balanceAt)
	bigBal := new(big.Float).Quo(fb, big.NewFloat(math.Pow10(18)))
	glog.Info("bigBal: ", bigBal)

	// get nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		glog.Error("pending nonce error: ", err)
	}
	// set value
	value := big.NewInt(10000000000000000) // in wei (0.01 eth)
	fvalue := new(big.Float)
	fvalue.SetInt(value)
	bigFBalance := new(big.Float).Quo(fvalue, big.NewFloat(math.Pow10(18)))
	glog.Info("bigFBalance: ", bigFBalance)

	// set gas limit
	gasLimit := uint64(21000)

	// get gas price
	suggestGasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		glog.Error("gas price error: ", err)
	}
	glog.Info("suggestGasPrice: ", suggestGasPrice)

	// get to address
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")

	// set data
	var data []byte

	// build tx
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: suggestGasPrice,
		Data:     data,
	})

	// get networkId
	networkID, err := client.NetworkID(context.Background())
	if err != nil {
		glog.Error("network id error: ", err)
	}
	glog.Info("networkID: ", networkID)

	// sign tx
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(networkID), privateKey)
	if err != nil {
		glog.Error("sign tx error: ", err)
	}

	_ = signTx

	// broadcast transaction
	//err = client.SendTransaction(context.Background(), signTx)
	//if err != nil {
	//	glog.Error(err)
	//}
	//
	//glog.Info("signTx sent: ", signTx.Hash().Hex())

	txHash := common.HexToHash("0x331ecee77fe0b77894bd3e5c600004c6d6adf22d5b52d4dade65fa41af129026")
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		glog.Error("transaction by hash error: ", err)
	}
	glog.Info("isPending: ", isPending)
	glog.Info("tx: ", tx)
}
