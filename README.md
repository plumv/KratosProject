Go-Kratos学习项目


# 项目结构

```text
  .
├── Dockerfile
├── LICENSE
├── Makefile
├── README.md
├── api // 下面维护了微服务使用的proto文件以及根据它们所生成的go文件
│ └── helloworld
│     └── v1
│         ├── error_reason.pb.go
│         ├── error_reason.proto
│         ├── error_reason.swagger.json
│         ├── greeter.pb.go
│         ├── greeter.proto
│         ├── greeter.swagger.json
│         ├── greeter_grpc.pb.go
│         └── greeter_http.pb.go
├── cmd  // 整个项目启动的入口文件
│ └── server
│     ├── main.go
│     ├── wire.go  // 我们使用wire来维护依赖注入
│     └── wire_gen.go
├── configs  // 这里通常维护一些本地调试用的样例配置文件
│ └── config.yaml
├── generate.go
├── go.mod
├── go.sum
├── internal  // 该服务所有不对外暴露的代码，通常的业务逻辑都在这下面，使用internal避免错误引用
│   ├── biz   // 业务逻辑的组装层，类似 DDD 的 domain 层，data 类似 DDD 的 repo，而 repo 接口在这里定义，使用依赖倒置的原则。
│ │ ├── README.md
│ │ ├── biz.go
│ │ └── greeter.go
│ ├── conf  // 内部使用的config的结构定义，使用proto格式生成
│ │ ├── conf.pb.go
│ │ └── conf.proto
│ ├── data  // 业务数据访问，包含 cache、db 等封装，实现了 biz 的 repo 接口。
│ │ ├── README.md
│ │ ├── data.go
│ │ └── greeter.go
│ ├── server  // http和grpc实例的创建和配置
│ │ ├── grpc.go
│ │ ├── http.go
│ │ └── server.go
│ └── service  // 实现了 api 定义的服务层，类似 DDD 的 application 层，处理 DTO 到 biz 领域实体的转换(DTO -> DO)，同时协同各类 biz 交互，但是不应处理复杂逻辑
│     ├── README.md
│     ├── greeter.go
│     └── service.go
└── third_party  // api 依赖的第三方proto
    ├── README.md
    ├── google
    │ └── api
    │     ├── annotations.proto
    │     ├── http.proto
    │     └── httpbody.proto
    └── validate
        ├── README.md
        └── validate.proto
```

特别说明：
biz层中定义了实体类，实体操作接口Repo，实体Usecase。
data层定义了“实体操作接口Repo”的实现。并支持通过data进行增强。
service层可以通过biz.实体Usecase来操作数据。
实体Usecase：提供数据操作逻辑。
实体操作接口Repo：提供操作数据的sql实现。

## 服务访问流程

RPC服务：浏览器 --> server/grpc.go --> service --> biz --> data
HTTP服务：浏览器 --> server/http.go --> service --> biz --> data


# proto文件说明

```prototext


```

# 创建项目

```shell
# 创建一个基本的项目
kratos new boss -r https://gitee.com/go-kratos/kratos-layout.git --nomod
# 删除无用的代码
rm -rf boss/LICENSE boss/openapi.yaml boss/api/helloworld
rm -rf boss/internal/biz/greeter.go boss/internal/biz/README.md boss/internal/biz/README.md   boss/api/helloword 
rm -rf boss/internal/data/greeter.go boss/internal/data/README.md boss/internal/service/greeter.go boss/internal/service/README.md
rm -rf boss/third_party/README.md boss/README.md
# 添加实体的proto文件
cd boss
kratos proto add api/base/base.proto
kratos proto client api/base/base.proto
kratos proto add api/user/user.proto
# 调整user.proto的内容，定义增删改查接口
kratos proto client api/user/user.proto
kratos proto server api/user/user.proto -t internal/service
# 调整下面文件
# boss/internal/service/service.go 将NewUserService加入为服务提供者
# boss/internal/server/grpc.go 将Service注入服务中进去
# boss/internal/server/http.go 将Service注入服务中进去
# 手写biz中的user.go文件和内容 并调整biz.go文件内容
# 手写data中的user.go文件和内容 并调整data.go文件内容
# 在boss/internal/service/service.go 实现通过biz.UserUsecase构造UserService

# 重新生成boss/cmd/boss/wire_gen.go文件
go get github.com/google/wire/cmd/wire@latest
go run github.com/google/wire/cmd/wire ./...

```



