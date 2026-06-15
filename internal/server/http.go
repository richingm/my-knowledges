package server

import (
	"encoding/json"
	"fmt"
	"io"
	fileV1 "my_knowledges/api/file/v1"
	v1 "my_knowledges/api/helloworld/v1"
	knowledgev1 "my_knowledges/api/knowledge/v1"
	userv1 "my_knowledges/api/user/v1"
	"my_knowledges/internal/conf"
	"my_knowledges/internal/service"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	httpTransport "github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, domain *service.DomainService, knowledge *service.KnowledgeService, article *service.ArticleService, file *service.FileService, user *service.UserService, logger log.Logger) *httpTransport.Server {
	var opts = []httpTransport.ServerOption{
		httpTransport.Middleware(
			recovery.Recovery(),
		),
		// 添加静态文件服务
		httpTransport.Filter(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 处理静态文件请求
				if r.URL.Path == "/" {
					w.Write([]byte("Welcome to Upload Service"))
					return
				}
				// 处理测试文件访问
				if r.URL.Path == "/upload_test.html" {
					// 提供测试文件
					http.ServeFile(w, r, "./upload_test.html")
					return
				}
				// 处理上传文件的静态访问
				if len(r.URL.Path) > 8 && r.URL.Path[:8] == "/uploads" {
					// 构建文件路径
					filePath := filepath.Join("./uploads", r.URL.Path[8:])
					// 提供静态文件
					http.ServeFile(w, r, filePath)
					return
				}
				// 处理 /files 路径的静态访问
				if len(r.URL.Path) > 7 && r.URL.Path[:7] == "/files/" {
					// 构建文件路径
					filePath := filepath.Join("./uploads", r.URL.Path[7:])
					// 提供静态文件
					http.ServeFile(w, r, filePath)
					return
				}
				// 处理文件上传请求
				if (r.URL.Path == "/upload/file/http" || r.URL.Path == "/upload/image/http" || r.URL.Path == "/api/v1/files/upload") && r.Method == http.MethodPost {
					// 对于multipart/form-data请求，我们需要直接处理
					// 这里我们将请求重定向到一个特殊的处理函数
					if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
						// 解析multipart表单
						if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB
							http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
							return
						}

						// 检查是否是文件上传请求
						if r.URL.Path == "/api/v1/files/upload" {
							// 直接处理文件上传请求
							file, header, err := r.FormFile("file")
							if err != nil {
								http.Error(w, "form file: "+err.Error(), http.StatusBadRequest)
								return
							}
							defer file.Close()

							filename := filepath.Base(header.Filename)
							now := time.Now()
							filename = fmt.Sprintf("%s-%s", now.Format("20060102150405"), filename)
							dst := "./uploads/" + filename

							// 创建上传目录
							if err := os.MkdirAll("./uploads", 0755); err != nil {
								http.Error(w, "create upload directory: "+err.Error(), http.StatusInternalServerError)
								return
							}

							// 创建目标文件
							out, err := os.Create(dst)
							if err != nil {
								http.Error(w, "create file: "+err.Error(), http.StatusInternalServerError)
								return
							}
							defer out.Close()

							// 复制文件内容
							if _, err := io.Copy(out, file); err != nil {
								http.Error(w, "copy file: "+err.Error(), http.StatusInternalServerError)
								return
							}

							// 生成文件URL
							fileURL := "http://localhost:8000/api/v1/files/" + filename

							// 返回上传结果
							reply := &fileV1.UploadReply{
								Name: filename,
								Url:  fileURL,
							}

							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(reply)
							return
						}
					}
				}
				// 处理 /api/v1/files 路径的静态访问
				if len(r.URL.Path) > 14 && r.URL.Path[:14] == "/api/v1/files/" {
					// 构建文件路径
					filePath := filepath.Join("./uploads", r.URL.Path[14:])
					// 提供静态文件
					http.ServeFile(w, r, filePath)
					return
				}
				// 其他请求交给下一处理器
				next.ServeHTTP(w, r)
			})
		}),
	}
	if c.Http.Network != "" {
		opts = append(opts, httpTransport.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, httpTransport.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, httpTransport.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := httpTransport.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	knowledgev1.RegisterDomainServiceHTTPServer(srv, domain)
	knowledgev1.RegisterKnowledgeServiceHTTPServer(srv, knowledge)
	knowledgev1.RegisterArticleServiceHTTPServer(srv, article)
	fileV1.RegisterFileServiceHTTPServer(srv, file)
	userv1.RegisterUserServiceHTTPServer(srv, user)
	return srv
}
