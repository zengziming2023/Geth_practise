package main

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/glog"
	"main/main/constant"
	"main/main/util"
	"main/token"
	"math/big"
	"strings"
)

type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

type LogApproval struct {
	TokenOwner common.Address
	Spender    common.Address
	Tokens     *big.Int
}

func main() {
	util.GlogInit()
	defer glog.Flush()

	client, err := ethclient.Dial(constant.SepoliaTest)
	if err != nil {
		glog.Fatal("error connecting to sepolia: ", err)
	}

	contractAddr := common.HexToAddress(constant.USDCContractAddr) // mtk address
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(6343663),
		ToBlock:   big.NewInt(6343674),
		Addresses: []common.Address{contractAddr},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		glog.Fatal("error filtering logs: ", err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(token.TokenMetaData.ABI))
	if err != nil {
		glog.Fatal("error parsing abi: ", err)
	}

	glog.Info("logs size = ", len(logs))
	for _, vlog := range logs {
		glog.Info("blockNumber: ", vlog.BlockNumber)
		glog.Info("index: ", vlog.Index)

		// 为了按某种日志类型进行过滤，我们需要弄清楚每个事件日志函数签名的 keccak256 哈希值
		logTransferSig := []byte("Transfer(address,address,uint256)")
		logApprovalSig := []byte("Approval(address,address,uint256)")
		logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
		logApprovalSigHash := crypto.Keccak256Hash(logApprovalSig)

		switch vlog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			glog.Info("find transfer ")
			logTransfer := LogTransfer{}
			err := contractAbi.UnpackIntoInterface(&logTransfer, "Transfer", vlog.Data)
			if err != nil {
				glog.Fatal("error unpacking transfer: ", err)
			}
			// 解包不会解析 indexed 事件类型，因为它们存储在 topics 下，所以对于那些我们必须单独解析
			logTransfer.From = common.HexToAddress(vlog.Topics[1].Hex())
			logTransfer.To = common.HexToAddress(vlog.Topics[2].Hex())
			glog.Info("logTransfer: ", logTransfer)
			glog.Info("tx hash: ", vlog.TxHash.Hex())
		case logApprovalSigHash.Hex():
			glog.Info("find approve ")
			logApproval := LogApproval{}
			err := contractAbi.UnpackIntoInterface(&logApproval, "Approval", vlog.Data)
			if err != nil {
				glog.Fatal("error unpacking approve: ", err)
			}
			logApproval.TokenOwner = common.HexToAddress(vlog.Topics[1].Hex())
			logApproval.Spender = common.HexToAddress(vlog.Topics[2].Hex())
			glog.Info("logApproval: ", logApproval)
		default:
			glog.Info("unknown log topic: ", vlog.Topics[0].Hex())
		}
	}
}
