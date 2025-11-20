package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
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
		field.Uint64("id").
			Immutable().
			Comment("主键"),
		field.Time("delete_at").Optional().Nillable().Comment("删除时间"),
		field.Uint64("created_by").Optional().Comment("创建人"),
		field.Uint64("updated_by").Optional().Comment("更新人"),
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
		index.Fields("delete_at"),
	}
}
