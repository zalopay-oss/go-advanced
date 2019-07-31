package main

import (
	"sync/atomic"
	"time"
)

func main() {
	var config atomic.Value // 保存当前配置信息

	// 初始化配置信息
	config.Store(loadConfig())

	// 启动一个后台线程, 加载更新后的配置信息
	go func() {
		for {
			time.Sleep(time.Second)
			config.Store(loadConfig())
		}
	}()

	// 用于处理请求的工作者线程始终采用最新的配置信息
	for i := 0; i < 10; i++ {
		go func() {
			for r := range requests() {
				c := config.Load()
				// ...
			}
		}()
	}
}
