package http

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

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
	// 从查询参数获取文件名，默认为index.md
	filename := r.URL.Query().Get("name")
	if filename == "" {
		filename = "index.md"
	}

	// 防止路径遍历攻击，确保文件在blogs目录下
	blogPath := filepath.Join("blogs", filename)
	if !strings.HasPrefix(blogPath, "blogs/") {
		http.Error(w, "无效的文件路径", http.StatusBadRequest)
		return
	}

	// 使用blogService获取博客内容
	contentStr, found, err := s.blogService.GetBlogContent(filename)
	if err != nil {
		http.Error(w, "无法读取博客文件: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	if !found {
		http.Error(w, "博客文件不存在", http.StatusNotFound)
		return
	}
	
	content := []byte(contentStr)

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
