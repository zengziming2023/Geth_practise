package main

import (
	"context"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/whisper/shhclient"
	"github.com/ethereum/go-ethereum/whisper/whisperv6"
	"github.com/golang/glog"
	"main/main/util"
)

func main() {
	// geth remove whisper from 1.10.0, pls noted.
	util.GlogInit()
	defer glog.Flush()

	glog.Info("whisper client start")

	// Whisper 是一种简单的基于点对点身份的消息传递系统，旨在成为下一去中心化的应用程序的构建块。 它旨在以相当的代价提供弹性和隐私。
	// 连接 Whisper 客户端
	// geth --ws  --shh // --rpc
	client, err := shhclient.Dial("ws://127.0.0.1:8546")
	if err != nil {
		glog.Fatal("whisper client dial error:", err)
	}
	_ = client
	glog.Info("whisper client success")

	keyID, err := client.NewKeyPair(context.Background())
	if err != nil {
		glog.Fatal("whisper client create key pair error:", err)
	}
	glog.Info("keyId: ", keyID)

	// get public key
	publicKey, err := client.PublicKey(context.Background(), keyID)
	if err != nil {
		glog.Fatal("whisper client get public key error:", err)
	}
	glog.Info("publicKey: ", hexutil.Encode(publicKey))

	//Payload 字节格式的消息内容
	//PublicKey 加密的公钥
	//TTL 消息的活跃时间
	//PowTime 做工证明的时间上限
	//PowTarget 做工证明的时间下限
	message := whisperv6.NewMessage{
		Payload:   []byte("hello"),
		PublicKey: publicKey,
		TTL:       60,
		PowTarget: 2.5,
		PowTime:   2,
	}

	messageHash, err := client.Post(context.Background(), message)
	if err != nil {
		glog.Fatal("whisper client post error:", err)
	}
	glog.Info("messageHash: ", messageHash)

	messages := make(chan *whisperv6.Message)
	// 消息的过滤标准
	criteria := whisperv6.Criteria{
		PrivateKeyID: keyID,
	}
	subscribeMessages, err := client.SubscribeMessages(context.Background(), criteria, messages)
	if err != nil {
		glog.Fatal("whisper client subscribe error:", err)
	}

	for {
		select {
		case msg := <-messages:
			glog.Info("whisper client receive message: ", msg)
		case err := <-subscribeMessages.Err():
			glog.Fatal("whisper client subscribe error:", err)
		}
	}
}
