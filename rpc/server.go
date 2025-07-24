package rpc

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "rpc-server/rpc-server/pb"
	"rpc-server/services"

	"google.golang.org/grpc"
)

// Server RPC服务器结构
type Server struct {
	pb.UnimplementedCalculatorServer
	blogService *services.BlogService
}

// NewServer 创建一个新的RPC服务器实例
func NewServer() *Server {
	return &Server{
		blogService: services.NewBlogService(),
	}
}

// Add 实现加法运算
func (s *Server) Add(ctx context.Context, in *pb.AddRequest) (*pb.AddResponse, error) {
	return &pb.AddResponse{Result: in.A + in.B}, nil
}

// GetBlog 实现获取博客内容
func (s *Server) GetBlog(ctx context.Context, in *pb.BlogRequest) (*pb.BlogResponse, error) {
	content, found, err := s.blogService.GetBlogContent(in.Name)
	if err != nil {
		return nil, err
	}
	
	return &pb.BlogResponse{
		Content: content,
		Found:   found,
	}, nil
}

// Start 启动RPC服务器
func (s *Server) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServer(grpcServer, s)
	log.Printf("gRPC服务器启动，监听端口: %d", port)

	return grpcServer.Serve(lis)
}