package hook

import (
	"boss/internal/data/ent"
	"context"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
)

var SoftDeleteHook = On(
	func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			// Skip soft-delete, means delete the entity permanently.
			if skip, _ := ctx.Value("softDeleteKey").(bool); skip {
				return next.Mutate(ctx, m)
			}
			mx, ok := m.(interface {
				SetOp(ent.Op)
				Client() *ent.Client
				SetDeleteTime(time.Time)
				WhereP(...func(*sql.Selector))
			})
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			mx.WhereP(
				sql.FieldIsNull("delete_at"),
			)
			mx.SetOp(ent.OpUpdate)
			mx.SetDeleteTime(time.Now())
			return mx.Client().Mutate(ctx, m)
		})
	},
	ent.OpDeleteOne|ent.OpDelete,
)
