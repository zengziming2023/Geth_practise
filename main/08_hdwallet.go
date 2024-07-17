package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"math/big"
)

func main() {
	mnemonic := "tag volcano eight thank tide danger coast health above argue embrace heavy"
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return
	}
	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return
	}
	fmt.Println(account.Address.Hex())

	path = hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/1")
	account, err = wallet.Derive(path, true)
	if err != nil {
		return
	}
	fmt.Println(account.Address.Hex())

	nonce := uint64(0)
	value := big.NewInt(1000000000000000000)
	toAddress := common.HexToAddress("0x0")
	gasLimit := uint64(21000)
	gasPrice := big.NewInt(21000000000)
	var data []byte

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &toAddress,
		Value:    value,
		Data:     data,
	})

	signTx, err := wallet.SignTx(account, tx, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	spew.Dump(signTx)
}
