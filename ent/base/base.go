package base

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// Base 定义基础的模型，以便实现一次定义，多处使用
// 由于mixin.Schema不具有传递性，因此只能在定义一遍
type Base struct {
	mixin.Schema
}

// Annotations of the User.
func (Base) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// 设定数据库表生成的时候加载注解
		entsql.WithComments(true),
	}
}

// Fields 字段定义
func (Base) Fields() []ent.Field {
	// 默认字段都是不为空的
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			// 业务侧兜底，防止 insert 时没带 id
			Default(uuid.New).
			Immutable().
			SchemaType(map[string]string{
				dialect.Postgres: "uuid DEFAULT gen_random_uuid()",
			}).Comment("主键"),
		field.Bool("is_delete").Default(false).Comment("是否删除:false表示不删除"),
		field.UUID("created_by", uuid.UUID{}).Optional().Comment("创建人"),
		field.UUID("updated_by", uuid.UUID{}).Optional().Comment("更新人"),
		field.Time("created_at").
			Immutable(). // 设置为不可变字段，即只在结构体创建的时候进行赋值
			Default(time.Now).
			Comment("创建时间"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).Comment("更新时间"),
	}
}

// Indexes 可按需加索引
func (Base) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("is_delete"),
	}
}
