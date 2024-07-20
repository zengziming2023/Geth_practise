package main

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/glog"
	"main/exchange"
	"main/main/constant"
	"main/main/util"
	"math/big"
	"strings"
)

// LogFill ...
type LogFill struct {
	Maker                  common.Address
	Taker                  common.Address
	FeeRecipient           common.Address
	MakerToken             common.Address
	TakerToken             common.Address
	FilledMakerTokenAmount *big.Int
	FilledTakerTokenAmount *big.Int
	PaidMakerFee           *big.Int
	PaidTakerFee           *big.Int
	Tokens                 [32]byte
	OrderHash              [32]byte
}

// LogCancel ...
type LogCancel struct {
	Maker                     common.Address
	FeeRecipient              common.Address
	MakerToken                common.Address
	TakerToken                common.Address
	CancelledMakerTokenAmount *big.Int
	CancelledTakerTokenAmount *big.Int
	Tokens                    [32]byte
	OrderHash                 [32]byte
}

type TokensPurchased struct {
	User          common.Address
	Amount        *big.Int
	Price         *big.Int
	PurchaseToken common.Address
	LockId        *big.Int
}

// LogError ...
type LogError struct {
	ErrorID   uint8
	OrderHash [32]byte
}

func main() {
	util.GlogInit()
	defer glog.Flush()

	//  solc --abi -o . Exchange.sol --overwrite
	// abigen --abi=Exchange.abi --pkg=exchange --out=../exchange/Exchange.go

	client, err := ethclient.Dial(constant.SepoliaTest)
	if err != nil {
		glog.Fatal("dial err:", err)
	}

	// 0x Protocol Exchange smart contract address
	contractAddr := common.HexToAddress("0x7940CEEB87B7a68cC9c2B29F5a596F95c9e7366F")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(6343814),
		ToBlock:   big.NewInt(6343881),
		Addresses: []common.Address{contractAddr},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		glog.Fatal("FilterLogs err:", err)
	}
	glog.Info("logs size: ", len(logs))

	contractAbi, err := abi.JSON(strings.NewReader(exchange.ExchangeMetaData.ABI))
	if err != nil {
		glog.Fatal("abi err:", err)
	}

	// NOTE: keccak256("LogFill(address,address,address,address,address,uint256,uint256,uint256,uint256,bytes32,bytes32)")
	logFillEvent := common.HexToHash("0d0b9391970d9a25552f37d436d2aae2925e2bfe1b2a923754bada030c498cb3")

	// NOTE: keccak256("LogCancel(address,address,address,address,uint256,uint256,bytes32,bytes32)")
	logCancelEvent := common.HexToHash("67d66f160bc93d925d05dae1794c90d2d6d6688b29b84ff069398a9b04587131")

	// NOTE: keccak256("LogError(uint8,bytes32)")
	logErrorEvent := common.HexToHash("36d86c59e00bd73dc19ba3adfe068e4b64ac7e92be35546adeddf1b956a87e90")

	// "TokensPurchased(address,uint256,uint256,address,uint256)"
	// crypto.Keccak256Hash([]byte("TokensPurchased(address,uint256,uint256,address,uint256)"))
	logTokensPurchasedEvent := common.HexToHash("0x0fbac4c1b53065ad309ad65b76cdd3048113d7def353fb1d2ab5d25928cef3fa")
	glog.Info(logTokensPurchasedEvent.Hex())

	for _, vlog := range logs {
		glog.Info("Log block num: ", vlog.BlockNumber)
		glog.Info("Log index: ", vlog.Index)
		glog.Info("Log topics[0]: ", vlog.Topics[0].Hex())

		switch vlog.Topics[0].Hex() {
		case logFillEvent.Hex():
			glog.Info("Log fill event: ", vlog)
		case logCancelEvent.Hex():
			glog.Info("Log cancel event: ", vlog)
		case logErrorEvent.Hex():
			glog.Info("Log error event: ", vlog)
		case logTokensPurchasedEvent.Hex():
			glog.Info("Log tokens purchase event: ", vlog)
			tokensPurchasedEvent := TokensPurchased{}
			err := contractAbi.UnpackIntoInterface(&tokensPurchasedEvent, "TokensPurchased", vlog.Data)
			if err != nil {
				glog.Fatal("unpack err:", err)
			}

			tokensPurchasedEvent.User = common.HexToAddress(vlog.Topics[1].Hex())
			tokensPurchasedEvent.Amount = new(big.Int).SetBytes(vlog.Topics[1][:])
			tokensPurchasedEvent.Price = new(big.Int).SetBytes(vlog.Topics[2][:])
			glog.Info("Log tx hex: ", vlog.TxHash.Hex())
			glog.Info("Tokens purchase event: ", tokensPurchasedEvent)
		default:
			glog.Info("Log unknown event: ", vlog)
		}

	}
}
