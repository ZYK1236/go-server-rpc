package services

import (
	"io/ioutil"
)

// BlogService 提供博客相关的服务
type BlogService struct{}

// NewBlogService 创建一个新的博客服务实例
func NewBlogService() *BlogService {
	return &BlogService{}
}

// GetBlogContent 根据文件名获取博客内容
func (bs *BlogService) GetBlogContent(name string) (string, bool, error) {
	// 构建文件路径
	blogPath := "blogs/" + name
	
	// 读取文件内容
	content, err := ioutil.ReadFile(blogPath)
	if err != nil {
		// 如果文件不存在，返回未找到
		return "", false, nil
	}
	
	// 返回文件内容
	return string(content), true, nil
}