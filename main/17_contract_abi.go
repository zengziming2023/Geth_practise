package main

func main() {
	// solc --abi -o . Store.sol 		// 生成 Store.abi
	// abigen --abi=Store.abi --pkg=store --out=../store/Store.go 	// store 文件夹需要先手动生成一下

	// 需要部署的话需要
	// 1.
	// solc --abi -o . Store.sol 		// 生成 Store.abi
	// 2.
	// solc --bin -o . Store.sol 		// 生成EVM字节码
	// abigen --bin=Store.bin --abi=Store.abi --pkg=store --out=../store/Store.go 	// 此时Store.go 文件会多出来 DeployStore 函数来部署合约
}
