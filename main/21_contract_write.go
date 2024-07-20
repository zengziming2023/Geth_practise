package main

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/glog"
	"main/main/constant"
	"main/main/util"
	"main/store"
	"math/big"
)

func main() {
	util.GlogInit()
	defer glog.Flush()

	client, err := ethclient.Dial(constant.SepoliaTest)
	if err != nil {
		glog.Fatal("dial:", err)
	}

	// get private key
	privateKey, err := crypto.HexToECDSA(constant.PrivateKey)
	if err != nil {
		glog.Fatal("private key:", err)
	}

	// get public key
	publicKey := privateKey.Public()
	// change to publicKeyECDSA
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		glog.Fatal("error casting public key to ECDSA")
	}

	// get from address
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// get nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		glog.Fatal("error calling PendingNonceAt:", err)
	}

	// get gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		glog.Fatal("error calling SuggestGasPrice:", err)
	}

	networkID, err := client.NetworkID(context.Background())
	if err != nil {
		glog.Fatal("error calling NetworkID:", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, networkID)
	if err != nil {
		glog.Fatal("error calling NewKeyedTransactorWithChainID:", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	contractAddr := common.HexToAddress(constant.StoreContractAddr)

	contractInstance, err := store.NewStore(contractAddr, client)
	if err != nil {
		glog.Fatal("error calling NewStore:", err)
	}

	key := [32]byte{}
	value := [32]byte{}
	copy(key[:], "foo")
	copy(value[:], "bar")

	//tx, err := contractInstance.SetItem(auth, key, value)
	//if err != nil {
	//	glog.Fatal("error calling SetItem:", err)
	//}
	//
	//glog.Info("tx:", tx.Hash().Hex())

	// get items
	result, err := contractInstance.Items(nil, key)
	if err != nil {
		glog.Fatal("error calling Items:", err)
	}
	glog.Info("result:", string(result[:]))
}
