package intercept

import (
	"context"

	"entgo.io/ent/dialect/sql"
)

// SoftDeleteInterceptor 软删除拦截器
var SoftDeleteInterceptor = TraverseFunc(func(ctx context.Context, q Query) error {
	// Skip soft-delete, means include soft-deleted entities.
	if skip, _ := ctx.Value("softDeleteKey").(bool); skip {
		return nil
	}
	// 默认添加查询条件delete_at为nil的条件
	q.WhereP(
		sql.FieldIsNull("delete_at"),
	)
	return nil
})
