package server

import (
	v1 "github.com/plumv/KratosProject/helloworld/api/helloworld/v1"
	"github.com/plumv/KratosProject/helloworld/internal/conf"
	"github.com/plumv/KratosProject/helloworld/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	// 创建http服务
	srv := http.NewServer(opts...)
	// 服务中注册Greeter服务：路由的绑定等逻辑是生成的。
	v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
