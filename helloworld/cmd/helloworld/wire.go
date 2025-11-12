//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/plumv/KratosProject/helloworld/internal/biz"
	"github.com/plumv/KratosProject/helloworld/internal/conf"
	"github.com/plumv/KratosProject/helloworld/internal/data"
	"github.com/plumv/KratosProject/helloworld/internal/server"
	"github.com/plumv/KratosProject/helloworld/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
// 定义使用wire的依赖方式
func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		// http和grpc实例的创建和配置
		server.ProviderSet,
		// 数据访问层实现
		data.ProviderSet,
		// 实体层和数据访问层接口
		biz.ProviderSet,
		// 业务服务层提供者
		service.ProviderSet,
		newApp))
}
