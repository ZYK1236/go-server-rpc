package main

import (
	"log"
	"sync"

	"rpc-server/http"
	"rpc-server/rpc"
)

func main() {
	// 创建等待组以等待所有goroutine完成
	var wg sync.WaitGroup

	// 启动RPC服务器
	wg.Add(1)
	go func() {
		defer wg.Done()
		rpcServer := rpc.NewServer()
		if err := rpcServer.Start(50051); err != nil {
			log.Fatalf("RPC服务器启动失败: %v", err)
		}
	}()

	// 启动HTTP服务器
	wg.Add(1)
	go func() {
		defer wg.Done()
		httpServer := http.NewServer()
		if err := httpServer.Start(8080); err != nil {
			log.Fatalf("HTTP服务器启动失败: %v", err)
		}
	}()

	// 等待所有服务器goroutine完成
	wg.Wait()
}
