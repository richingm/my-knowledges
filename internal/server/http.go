package server

import (
	"context"
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
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	httpTransport "github.com/go-kratos/kratos/v2/transport/http"
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

// jwtMiddleware 自定义 JWT 中间件，跳过登录注册接口
func jwtMiddleware(secret string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				if httpTr, ok := tr.(*httpTransport.Transport); ok {
					path := httpTr.Request().URL.Path
					if path == "/api/v1/users/login" || path == "/api/v1/users/register" {
						return handler(ctx, req)
					}
				}
			}
			return jwt.Server(func(token *jwtv5.Token) (interface{}, error) {
				return []byte(secret), nil
			})(handler)(ctx, req)
		}
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, data *conf.Data, greeter *service.GreeterService, domain *service.DomainService, knowledge *service.KnowledgeService, article *service.ArticleService, file *service.FileService, user *service.UserService, logger log.Logger) *httpTransport.Server {
	var opts = []httpTransport.ServerOption{
		httpTransport.Middleware(
			recovery.Recovery(),
			jwtMiddleware(data.Jwt.Secret),
		),
		httpTransport.Filter(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/" {
					http.ServeFile(w, r, "./web/dist/index.html")
					return
				}
				if r.URL.Path == "/upload_test.html" {
					http.ServeFile(w, r, "./upload_test.html")
					return
				}
				if len(r.URL.Path) > 8 && r.URL.Path[:8] == "/uploads" {
					filePath := filepath.Join("./uploads", r.URL.Path[8:])
					http.ServeFile(w, r, filePath)
					return
				}
				if len(r.URL.Path) > 7 && r.URL.Path[:7] == "/files/" {
					filePath := filepath.Join("./uploads", r.URL.Path[7:])
					http.ServeFile(w, r, filePath)
					return
				}
				if len(r.URL.Path) > 7 && r.URL.Path[:7] == "/assets" {
					filePath := filepath.Join("./web/dist", r.URL.Path)
					http.ServeFile(w, r, filePath)
					return
				}
				if len(r.URL.Path) > 14 && r.URL.Path[:14] == "/api/v1/files/" {
					filePath := filepath.Join("./uploads", r.URL.Path[14:])
					http.ServeFile(w, r, filePath)
					return
				}
				if (r.URL.Path == "/upload/file/http" || r.URL.Path == "/upload/image/http" || r.URL.Path == "/api/v1/files/upload") && r.Method == http.MethodPost {
					if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
						if err := r.ParseMultipartForm(10 << 20); err != nil {
							http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
							return
						}
						if r.URL.Path == "/api/v1/files/upload" {
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
							if err := os.MkdirAll("./uploads", 0755); err != nil {
								http.Error(w, "create upload directory: "+err.Error(), http.StatusInternalServerError)
								return
							}
							out, err := os.Create(dst)
							if err != nil {
								http.Error(w, "create file: "+err.Error(), http.StatusInternalServerError)
								return
							}
							defer out.Close()
							if _, err := io.Copy(out, file); err != nil {
								http.Error(w, "copy file: "+err.Error(), http.StatusInternalServerError)
								return
							}
							fileURL := "http://39.105.17.55:8000/api/v1/files/" + filename
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
