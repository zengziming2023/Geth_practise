package main

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/glog"
	"main/main/constant"
	"main/main/util"
	"main/store"
	"math/big"
	"strings"
)

func main() {
	util.GlogInit()
	defer glog.Flush()

	client, err := ethclient.Dial(constant.SepoliaTestRPC)
	if err != nil {
		glog.Fatal("dial eth err:", err)
	}

	contractAddr := common.HexToAddress(constant.StoreContractAddr)

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(6338418),
		ToBlock:   big.NewInt(6339631),
		Addresses: []common.Address{contractAddr},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		glog.Fatal("FilterLogs err:", err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(store.StoreMetaData.ABI))
	if err != nil {
		glog.Fatal("abi err:", err)
	}
	glog.Info("contract abi:", contractAbi)

	for _, vLog := range logs {
		glog.Info("block hash: ", vLog.BlockHash.Hex())
		glog.Info("block number: ", vLog.BlockNumber)
		glog.Info("tx hash: ", vLog.TxHash.Hex())
		glog.Info("data: ", hex.EncodeToString(vLog.Data))

		event := struct {
			Key   [32]byte // 必须是首字母大写，这样才是Export的，反射才可以应用上
			Value [32]byte
		}{}

		err = contractAbi.UnpackIntoInterface(&event, "ItemSet", vLog.Data)
		if err != nil {
			glog.Fatal("abi.Unpack err:", err)
		}

		glog.Info("Key: ", string(event.Key[:]))
		glog.Info("Value: ", string(event.Value[:]))

		var topics [4]string
		for i := range vLog.Topics {
			topics[i] = vLog.Topics[i].Hex()
		}
		glog.Info("topics[0]: ", topics[0])

	}

	eventSignature := []byte("ItemSet(bytes32,bytes32)")
	hash := crypto.Keccak256Hash(eventSignature)
	glog.Info("hash: ", hash.Hex()) // should equal with topics[0]

}
