package main

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/glog"
	"golang.org/x/crypto/sha3"
	"main/main/constant"
	"main/main/util"
	"math"
	"math/big"
)

func main() {
	util.GlogInit()
	defer glog.Flush()

	glog.Info("Starting token Transfer...")
	defer glog.Info("Stopping token Transfer...")
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

	glog.Info("nonce: ", nonce)

	// token transfer, eth value is zero
	value := big.NewInt(0)

	// get suggest gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		glog.Error("gas price: ", err)
	}

	// transfer to who
	toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	// token address
	tokenContractAddr := common.HexToAddress("0x28b149020d2152179873ec60bed6bf7cd705775d")

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	glog.Info("methodID: ", hexutil.Encode(methodID)) // 0xa9059cbb

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	glog.Info("paddedAddress: ", hexutil.Encode(paddedAddress)) // 0x0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d

	amount := new(big.Int)
	amount.SetString("1000000000000000000000", 10)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	glog.Info("paddedAmount: ", hexutil.Encode(paddedAmount)) //0x00000000000000000000000000000000000000000000003635c9adc5dea00000

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	// gas fee
	estimateGas, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &tokenContractAddr,
		Data: data,
	})
	if err != nil {
		glog.Error("estimate gas error: ", err)
	}
	glog.Info("estimate gas: ", estimateGas)

	// build tx
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &tokenContractAddr,
		Value:    value,
		GasPrice: gasPrice,
		Gas:      estimateGas,
		Data:     data,
	})

	// get network id
	networkID, err := client.NetworkID(context.Background())
	if err != nil {
		glog.Error("client network id error: ", err)
	}
	glog.Info("networkID: ", networkID)

	// sign tx
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(networkID), privateKey)
	if err != nil {
		glog.Error("sign tx error: ", err)
	}
	glog.Info("signTx: ", signTx.Hash().Hex()) // 0x329a8bc52075b09811cdbfcc44fb41a0f4cbeb78a9715a0ef26e759feae5e3d4

	// broadcast tx
	//err = client.SendTransaction(context.Background(), signTx)
	//if err != nil {
	//	glog.Error("send tx error: ", err)
	//}
}
