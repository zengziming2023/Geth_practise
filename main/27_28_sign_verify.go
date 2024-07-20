package main

import (
	"bytes"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/glog"
	"main/main/constant"
	"main/main/util"
)

func main() {
	util.GlogInit()
	defer glog.Flush()

	privateKey, err := crypto.HexToECDSA(constant.PrivateKey)
	if err != nil {
		glog.Fatal("hex to ECDSA:", err)
	}
	data := []byte("hello")
	hash := crypto.Keccak256Hash(data)
	glog.Info("hash:", hash.Hex())

	sign, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		glog.Fatal("sign:", err)
	}

	glog.Info("sign: ", hexutil.Encode(sign))

	// 有 3 件事来验证签名：签名，原始数据的哈希以及签名者的公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		glog.Fatal("error casting public key to ECDSA")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	sigPublicKeyBytes, err := crypto.Ecrecover(hash.Bytes(), sign)

	// 须将签名的公钥与期望的公钥进行比较
	match := bytes.Equal(publicKeyBytes, sigPublicKeyBytes)
	glog.Info("match:", match)

	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), sign)
	if err != nil {
		return
	}
	sigPublicKeyBytes = crypto.FromECDSAPub(sigPublicKeyECDSA)
	match = bytes.Equal(publicKeyBytes, sigPublicKeyBytes)
	glog.Info("match:", match)

	signatureNoRecoverID := sign[:len(sign)-1] //remove recovery id
	verified := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
	glog.Info("verified:", verified)
}
