package main

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

	client, err := ethclient.Dial(constant.SepoliaTest)
	if err != nil {
		glog.Fatal(err)
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
	glog.Info("fromAddress: ", fromAddress)

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
	glog.Info("nonce: ", nonce)

	networkID, err := client.NetworkID(context.Background())
	if err != nil {
		glog.Error("error getting network id: ", err)
	}
	glog.Info("networkID: ", networkID)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		glog.Fatal("error getting gas price: ", err)
	}

	// deploy contract
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, networkID)
	if err != nil {
		glog.Fatal("new keyed transactor error: ", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	//input := "1.0"
	//contractAddr, tx, instance, err := store.DeployStore(auth, client, input)
	//if err != nil {
	//	glog.Fatal("DeployStore error: ", err)
	//}
	//glog.Info("contractAddr: ", contractAddr.Hex()) 	// 0x2D526f19d8bBdAf09BFf4de187D2ad3822633b53
	//glog.Info("tx: ", tx.Hash().Hex()) 		// 0x0eac93bff346ec53444269b5fcdd272f56228385e9a22152f1f3a811cde72e4d
	//_ = instance
}
