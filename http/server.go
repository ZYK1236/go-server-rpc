package http

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"

	"rpc-server/services"
)

// Server HTTP服务器结构
type Server struct {
	blogService *services.BlogService
}

// NewServer 创建一个新的HTTP服务器实例
func NewServer() *Server {
	return &Server{
		blogService: services.NewBlogService(),
	}
}

// blogHandler 处理博客请求
func (s *Server) blogHandler(w http.ResponseWriter, r *http.Request) {
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
}

// Start 启动HTTP服务器
func (s *Server) Start(port int) error {
	// 注册处理函数
	http.HandleFunc("/blogs", s.blogHandler)
	
	log.Printf("HTTP服务器启动，监听端口: %d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}