Go-Kratos学习项目


# 项目结构

```text
  .
├── README.md // 整个项目的搭建文档
├── api // 下面维护了所有微服务使用的proto文件以及根据它们所生成的go文件
│ ├── base // 所有微服务基础的对外的请求响应结构
│ └── boss // boss微服务的proto文件
├── ent  // 下面维护了所有微服务使用的表的schema
│ ├── base // 所有微服务基础表结构
│ └── boss // boss微服务的表结构
├── boss // boss微服务
│ ├── cmd // 启动入口
│ ├── config // 配置信息
│ ├── go.mod  // boss服务的mod文件
│ └── internal // 该服务所有不对外暴露的代码
│   ├── biz // 业务逻辑的组装层
│   ├── conf // 配置文件对应的结构体
│   ├── data // 业务数据访问，包含 cache、db 等封装，实现了 biz 的 repo 接口。
│      └── ent // 采用ent框架实现业务数据的ORM
│   ├── server // http和grpc实例的创建和配置
│   └── service // 实现了 api 定义的服务层
├── third_party // api 依赖的第三方proto
├── go.mod // 整体项目的mod文件
└── go.work // 整个工作空间的work文件
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


## proto文件创建以及生成

```shell

# 添加实体的proto文件
mkdir api
cd api
go mod init
```

### 创建常用的请求响应结构

```shell
cd api
kratos proto add base/base.proto
kratos proto client base/base.proto
```

```protobuf
syntax = "proto3";

package base;

option go_package = "api/base;base";

message SortReq{
  string order  = 1;
  string sort  = 2;
}

message PageReq{
  int32 page_size = 1;
  int32 page = 2;
}

message ResultResp{
  int32 code  = 1;
  string msg = 2;
}

message PageResp{
  int64 total  = 1;
}

```

### 创建boos服务的user

```shell
cd api
kratos proto add boss/user.proto
kratos proto client boss/user.proto
```

```protobuf
syntax = "proto3";

package boss;
import "base/base.proto";

option go_package = "api/boss;boss";

message UserReply{
  int64  id = 1;
  string name = 2;
  int32 age = 3;
  string email = 4;
}

message CreateUserRequest {
  string name = 1;
  int32 age = 2;
  string email = 3;
  string password = 4;
}
message CreateUserReply {
  base.ResultResp r = 1;
  UserReply data = 2;
}

message UpdateUserRequest {
  int64  id = 1;
  string name = 2;
  int32 age = 3;
  string email = 4;
}
message UpdateUserReply {
  base.ResultResp r = 1;
  UserReply data = 2;
}

message DeleteUserRequest {
  int64  id = 1;
}
message DeleteUserReply {
  base.ResultResp r = 1;
}

message GetUserRequest {
  int64  id = 1;

}
message GetUserReply {
  base.ResultResp r = 1;
  UserReply data = 2;
}

message QueryUserRequest{
  string name = 1;
  int32 age = 2;
}

message ListUserRequest {
  base.SortReq sort = 1;
  QueryUserRequest query = 2;
}
message ListUserReply {
  base.ResultResp r = 1;
  repeated UserReply data = 2;
}

message PageUserRequest {
  base.PageReq page = 1;
  base.SortReq sort = 2;
  QueryUserRequest query = 3;
}

message PageUserReply{
  base.ResultResp r = 1;
  base.PageResp page = 2;
  repeated UserReply data = 3;
}
```

### 创建boos服务的role

```shell
cd api
kratos proto add boss/role.proto
kratos proto client boss/role.proto
```

```protobuf
syntax = "proto3";

package boss;
import "base/base.proto";

option go_package = "api/boss;boss";

message RoleReply{
  int64  id = 1;
  string name = 2;
  string desc = 3;
}
message CreateRoleRequest {
  string name = 1;
  string desc = 2;
}
message UpdateRoleRequest {
  int64  id = 1;
  string name = 2;
  string desc = 3;
}


message CreateRoleReply {
  base.ResultResp r = 1;
  RoleReply data = 2;
}
message UpdateRoleReply {
  base.ResultResp r = 1;
  RoleReply data = 2;
}

message DeleteRoleRequest {
  int64  id = 1;
}
message DeleteRoleReply {
  base.ResultResp r = 1;
}

message GetRoleRequest {
  int64  id = 1;

}
message GetRoleReply {
  base.ResultResp r = 1;
  RoleReply data = 2;
}

message QueryRoleRequest{
  string name = 1;
  int32 age = 2;
}

message ListRoleRequest {
  base.SortReq sort = 1;
  QueryRoleRequest query = 2;
}
message ListRoleReply {
  base.ResultResp r = 1;
  repeated RoleReply data = 2;
}

message PageRoleRequest {
  base.PageReq page = 1;
  base.SortReq sort = 2;
  QueryRoleRequest query = 3;
}

message PageRoleReply{
  base.ResultResp r = 1;
  base.PageResp page = 2;
  repeated RoleReply data = 3;
}
```

### 创建boos服务的boss

```shell
cd api
kratos proto add boss/boss.proto
kratos proto client boss/boss.proto
mv openapi.yaml boss/boss_openapi.yaml
```

```protobuf
syntax = "proto3";

package boss;

import "google/api/annotations.proto";
import "boss/user.proto";
import "boss/role.proto";

option go_package = "api/boss;boss";
option java_multiple_files = true;
option java_package = "boss";

service User {
  rpc CreateUser (CreateUserRequest) returns (CreateUserReply){
    option (google.api.http) = {
      post: "/api/user",
      body: "*"
    };
  };
  rpc UpdateUser (UpdateUserRequest) returns (UpdateUserReply){
    option (google.api.http) = {
      put: "/api/user/:id",
      body: "*"
    };
  };
  rpc DeleteUser (DeleteUserRequest) returns (DeleteUserReply){
    option (google.api.http) = {
      delete: "/api/user/:id"
    };
  };
  rpc GetUser (GetUserRequest) returns (GetUserReply){
    option (google.api.http) = {
      get: "/api/user/:id"
    };
  };
  rpc ListUser (ListUserRequest) returns (ListUserReply){
    option (google.api.http) = {
      post: "/api/user/list",
      body: "*"
    };
  };
  rpc PageUser (PageUserRequest) returns (PageUserReply){
    option (google.api.http) = {
      post: "/api/user/page",
      body: "*"
    };
  };
}

service Role {
  rpc CreateRole (CreateRoleRequest) returns (CreateRoleReply){
    option (google.api.http) = {
      post: "/api/role",
      body: "*"
    };
  };
  rpc UpdateRole (UpdateRoleRequest) returns (UpdateRoleReply){
    option (google.api.http) = {
      put: "/api/role/:id",
      body: "*"
    };
  };
  rpc DeleteRole (DeleteRoleRequest) returns (DeleteRoleReply){
    option (google.api.http) = {
      delete: "/api/role/:id"
    };
  };
  rpc GetRole (GetRoleRequest) returns (GetRoleReply){
    option (google.api.http) = {
      get: "/api/role/:id"
    };
  };
  rpc ListRole (ListRoleRequest) returns (ListRoleReply){
    option (google.api.http) = {
      post: "/api/role/list",
      body: "*"
    };
  };
  rpc PageRole (PageRoleRequest) returns (PageRoleReply){
    option (google.api.http) = {
      post: "/api/role/page",
      body: "*"
    };
  };
}
```

## 创建boss服务

### 生成项目框架

```shell
# 创建一个基本的项目
kratos new boss -r https://gitee.com/go-kratos/kratos-layout.git
# 删除无用的代码
rm -rf boss/LICENSE boss/openapi.yaml boss/api/helloworld
rm -rf boss/internal/biz/greeter.go boss/internal/biz/README.md boss/internal/biz/README.md   boss/api/helloword 
rm -rf boss/internal/data/greeter.go boss/internal/data/README.md boss/internal/service/greeter.go boss/internal/service/README.md
rm -rf boss/third_party/README.md boss/README.md
rm -rf boss/api boss/third_party
```
### 调整Biz层

创建User.go代码
```go
package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// User is a User model.
type User struct {
	Hello string
}

// UserRepo is a Greater repo.
type UserRepo interface {
	Save(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	FindByID(context.Context, int64) (*User, error)
	ListByHello(context.Context, string) ([]*User, error)
	ListAll(context.Context) ([]*User, error)
}

// UserUsecase is a User usecase.
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// NewUserUsecase new a User usecase.
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateUser creates a User, and returns the new User.
func (uc *UserUsecase) CreateUser(ctx context.Context, g *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("CreateUser: %v", g.Hello)
	return uc.repo.Save(ctx, g)
}

```

创建Role.go代码
```go
package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// Role is a Role model.
type Role struct {
	Hello string
}

// RoleRepo is a Greater repo.
type RoleRepo interface {
	Save(context.Context, *Role) (*Role, error)
	Update(context.Context, *Role) (*Role, error)
	FindByID(context.Context, int64) (*Role, error)
	ListByHello(context.Context, string) ([]*Role, error)
	ListAll(context.Context) ([]*Role, error)
}

// RoleUsecase is a Role usecase.
type RoleUsecase struct {
	repo RoleRepo
	log  *log.Helper
}

// NewRoleUsecase new a Role usecase.
func NewRoleUsecase(repo RoleRepo, logger log.Logger) *RoleUsecase {
	return &RoleUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateRole creates a Role, and returns the new Role.
func (uc *RoleUsecase) CreateRole(ctx context.Context, g *Role) (*Role, error) {
	uc.log.WithContext(ctx).Infof("CreateRole: %v", g.Hello)
	return uc.repo.Save(ctx, g)
}

```


将Usecase的创建方法通过wire暴露调整biz.go文件
```go
package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewRoleUsecase, NewUserUsecase)
```

### 调整Data层

创建User.go代码

```go
package data

import (
	"boss/internal/biz"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) Save(ctx context.Context, g *biz.User) (*biz.User, error) {
	return g, nil
}

func (r *userRepo) Update(ctx context.Context, g *biz.User) (*biz.User, error) {
	return g, nil
}

func (r *userRepo) FindByID(context.Context, int64) (*biz.User, error) {
	return nil, nil
}

func (r *userRepo) ListByHello(context.Context, string) ([]*biz.User, error) {
	return nil, nil
}

func (r *userRepo) ListAll(context.Context) ([]*biz.User, error) {
	return nil, nil
}

```

创建Role.go代码

```go
package data

import (
	"boss/internal/biz"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type roleRepo struct {
	data *Data
	log  *log.Helper
}

// NewRoleRepo .
func NewRoleRepo(data *Data, logger log.Logger) biz.RoleRepo {
	return &roleRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *roleRepo) Save(ctx context.Context, g *biz.Role) (*biz.Role, error) {
	return g, nil
}

func (r *roleRepo) Update(ctx context.Context, g *biz.Role) (*biz.Role, error) {
	return g, nil
}

func (r *roleRepo) FindByID(context.Context, int64) (*biz.Role, error) {
	return nil, nil
}

func (r *roleRepo) ListByHello(context.Context, string) ([]*biz.Role, error) {
	return nil, nil
}

func (r *roleRepo) ListAll(context.Context) ([]*biz.Role, error) {
	return nil, nil
}

```

将Rep的创建方法通过wire暴露 调整data.go文件

```go
package data

import (
	"boss/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo,NewRoleRepo)

// Data .
type Data struct {
	// TODO wrapped database client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{}, cleanup, nil
}

```

### 调整service层代码

生成service代码

```shell
kratos proto server api/boss/boss.proto -t boss/internal/service
```

调整user.go文件将biz的UserUsecase注入到UserService中

```go
type UserService struct {
	pb.UnimplementedUserServer
	uc *biz.UserUsecase
}

func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{
		uc: uc,
	}
}
```

调整role.go文件将biz的RoleUsecase注入到RoleService中

```go
type RoleService struct {
pb.UnimplementedRoleServer
uc *biz.RoleUsecase
}

func NewRoleService(uc *biz.RoleUsecase) *RoleService {
return &RoleService{
uc: uc,
}
}
```

将service代码添加到service.go中通过wire暴露

```go
package service

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewUserService,NewRoleService)

```




### 调整server层代码

将service层的服务注入到grpc服务和http服务中

```go
package server

import (
	"api/boss"
	"boss/internal/conf"
	"boss/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, 
	userService *service.UserService, 
	roleService *service.RoleService, 
	logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
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
	srv := grpc.NewServer(opts...)
	boss.RegisterUserServer(srv, userService)
	boss.RegisterRoleServer(srv, roleService)
	return srv
}

```

```go
package server

import (
	"api/boss"
	"boss/internal/conf"
	"boss/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, 
	userService *service.UserService,
	roleService *service.RoleService,
	logger log.Logger) *http.Server {
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
	srv := http.NewServer(opts...)
	boss.RegisterUserHTTPServer(srv, userService)
	boss.RegisterRoleHTTPServer(srv, roleService)
	return srv
}

```

### 调整cmd下的代码

删除wire_gen.go文件，并执行以下命令，重新生成wire_gen.go文件

```shell
cd boss
rm -rf cmd/boss/wire_gen.go
go run github.com/google/wire/cmd/wire ./...
```

## 运行boss服务

### 命令行启动方式

```shell
cd boss
kratos run
```

### IDEA启动方式

调整config文件路径地址为当前项目下的路径，而不是main函数相对路径。

```go
func init() {
	flag.StringVar(&flagconf, "conf", "configs", "config path, eg: -conf config.yaml")
}
```

双击cmd/boss/main.go文件中的main函数即可


# 完善项目代码

## 完善数据库ORM框架:ent

生成用户和角色的实体表

    ent generate --target <target dirpath> <template dirpath>

```shell
go get entgo.io/ent
go run entgo.io/ent/cmd/ent new User Role
# 参考./ent/schema 生成到 ./boss/internal/data/ent
go run entgo.io/ent/cmd/ent generate --target  ./boss/internal/data/ent ./ent/boss

go run entgo.io/ent/cmd/ent generate --target  ./ent/base ./ent/base/schema

```



#  截断


```shell
kratos proto add base/base.proto
kratos proto client api/base/base.proto
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

替换命令  kratos proto client boss/boss.proto --还是不替换好。
protoc --proto_path=. --proto_path=../third_party --openapi_out=fq_schema_naming=true,default_response=false:./boss boss/boss.proto
