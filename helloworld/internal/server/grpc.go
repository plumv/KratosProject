package server

import (
	v1 "github.com/plumv/KratosProject/helloworld/api/helloworld/v1"
	"github.com/plumv/KratosProject/helloworld/internal/conf"
	"github.com/plumv/KratosProject/helloworld/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			// 恢复中间件
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	// 构建grpc服务
	srv := grpc.NewServer(opts...)
	// 服务中注册Greeter服务
	v1.RegisterGreeterServer(srv, greeter)
	return srv
}
