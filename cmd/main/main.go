package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"io/ioutil"

	pb "rpc-server/rpc-server/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCalculatorServer
}

func (s *server) Add(ctx context.Context, in *pb.AddRequest) (*pb.AddResponse, error) {
	return &pb.AddResponse{Result: in.A + in.B}, nil
}

func (s *server) GetBlog(ctx context.Context, in *pb.BlogRequest) (*pb.BlogResponse, error) {
	// 构建文件路径
	blogPath := "blogs/" + in.Name
	
	// 读取文件内容
	content, err := ioutil.ReadFile(blogPath)
	if err != nil {
		// 如果文件不存在，返回未找到
		return &pb.BlogResponse{Found: false}, nil
	}
	
	// 返回文件内容
	return &pb.BlogResponse{
		Content: string(content),
		Found:   true,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCalculatorServer(s, &server{})
	// 启动HTTP服务器来提供博客内容
	go func() {
		http.HandleFunc("/blogs", func(w http.ResponseWriter, r *http.Request) {
			// 使用相对路径读取blogs/index.md文件
			blogPath := "blogs/index.md"
			
			// 读取文件内容
			content, err := ioutil.ReadFile(blogPath)
			if err != nil {
				http.Error(w, "无法读取博客文件: " + err.Error(), http.StatusInternalServerError)
				return
			}
			
				// 设置CORS头
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			
			// 处理预检请求
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			// 设置响应头并返回内容
			w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write(content)
		})
		
		log.Printf("HTTP服务器启动，监听端口: %d", 8080)
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("HTTP服务启动失败: %v", err)
		}
	}()

	log.Printf("gRPC服务器启动，监听端口: %d", 50051)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("gRPC服务失败: %v", err)
	}
}
