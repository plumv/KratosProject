package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// User 用户表
type User struct {
	ent.Schema
}

// Mixin 混入的公共字段
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Base{},
	}
}

// Annotations 额外的定义
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// 设定数据库表生成的时候加载注解
		entsql.Table("boss_user"),
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").Comment("用户名"),
		field.String("password").Comment("密码"),
		field.Int32("age").Comment("年龄"),
	}
}
