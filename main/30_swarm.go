package main

func main() {
	// Swarm 是以太坊的去中心化和分布式的存储解决方案，与 IPFS 类似。

	// 	•	go get -d：只下载源码，不编译和安装。
	//	•	go get -u：更新包及其所有依赖包到最新版本。

	//go get -u github.com/ethereum/go-ethereum
	//go get -u github.com/ethereum/go-ethereum/cmd/geth
	//go get -u github.com/ethereum/go-ethereum/cmd/swarm 	// swarm 不再是go-ethereum 下了

	//go get -u github.com/ethersphere/swarm
	//geth account new
	// password: hele
	//Your new key was generated
	//
	//Public address of the key:   0x44b7B3293C829cF2b440a63Bc14DE902fEa53733
	//Path of the secret key file: /Users/ziming.zeng/Library/Ethereum/keystore/UTC--2024-07-21T03-26-50.913836000Z--44b7b3293c829cf2b440a63bc14de902fea53733
	//
	//- You can share your public address with anyone. Others need it to interact with you.
	//- You must NEVER share the secret key with anyone! The key controls access to your funds!
	//- You must BACKUP your key file! Without the key, it's impossible to access account funds!
	//- You must REMEMBER your password! Without the password, it's impossible to decrypt the key!

	// 启动swarm 服务
	//export BZZKEY=44b7B3293C829cF2b440a63Bc14DE902fEa53733
	//swarm --bzzaccount $BZZKEY
}
