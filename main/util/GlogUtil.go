package util

import (
	"flag"
)

func GlogInit() {
	flag.Parse()
	flag.Set("logtostderr", "true")
	flag.Set("v", "4") // 设置日志级别
}
