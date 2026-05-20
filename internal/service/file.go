package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	v1 "my_knowledges/api/file/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// FileService 文件服务
type FileService struct {
	v1.UnimplementedFileServiceServer
	log       *log.Helper
	uploadDir string
}

// NewFileService 创建文件服务
func NewFileService(logger log.Logger) *FileService {
	uploadDir := "./uploads"
	// 创建上传目录
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		logger.Log(log.LevelError, "err", err)
	}
	return &FileService{
		log:       log.NewHelper(logger),
		uploadDir: uploadDir,
	}
}

// Upload 上传文件
func (s *FileService) Upload(ctx context.Context, in *v1.UploadRequest) (*v1.UploadReply, error) {
	s.log.Info("Upload file")
	// 文件上传逻辑已在 HTTP 服务器的过滤器中实现
	return &v1.UploadReply{
		Name: "test.txt",
		Url:  "/files/test.txt",
	}, nil
}

// View 查看文件
func (s *FileService) View(ctx context.Context, in *v1.ViewRequest) (*v1.ViewReply, error) {
	s.log.Info("View file: %s", in.Filename)

	// 构建文件路径
	filepath := filepath.Join(s.uploadDir, in.Filename)

	// 读取文件内容
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 确定文件的 MIME 类型
	contentType := http.DetectContentType(content)

	return &v1.ViewReply{
		Content:     content,
		ContentType: contentType,
	}, nil
}

// Download 下载文件
func (s *FileService) Download(ctx context.Context, in *v1.ViewRequest) (*v1.ViewReply, error) {
	s.log.Info("Download file: %s", in.Filename)

	// 构建文件路径
	filepath := filepath.Join(s.uploadDir, in.Filename)

	// 读取文件内容
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 确定文件的 MIME 类型
	contentType := http.DetectContentType(content)

	return &v1.ViewReply{
		Content:     content,
		ContentType: contentType,
	}, nil
}
