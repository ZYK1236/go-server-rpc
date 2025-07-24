package main

import (
	"context"
	"log"
	"time"

	pb "rpc-server/rpc-server/proto"

	"google.golang.org/grpc"
)

func main() {
	// 连接到gRPC服务器
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("无法连接到gRPC服务器: %v", err)
	}
	defer conn.Close()

	// 创建gRPC客户端
	client := pb.NewCalculatorClient(conn)

	// 调用GetBlog方法
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 请求获取index.md文件
	blogRequest := &pb.BlogRequest{Name: "index.md"}
	r, err := client.GetBlog(ctx, blogRequest)
	if err != nil {
		log.Fatalf("无法获取博客: %v", err)
	}

	if r.Found {
		log.Printf("成功获取博客内容:\n%s", r.Content)
	} else {
		log.Printf("未找到指定的博客文件")
	}
}