package biz

import (
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewRoleUsecase, NewUserUsecase)

// Page 分页查询
type Page struct {
	Page   *int32
	Limit  *int32
	Orders *[]*Order
}

// Order 排序
type Order struct {
	Field *string
	Desc  *bool
}
