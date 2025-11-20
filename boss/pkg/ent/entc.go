//go:build ignore

package main

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/go-kratos/kratos/v2/log"
)

func main() {
	schemaPath := "./schema"
	genConfig := &gen.Config{
		//IDType: &field.TypeInfo{Type: field.TypeInt},
		Features: []gen.Feature{
			gen.FeatureEntQL,
			gen.FeaturePrivacy,
			gen.FeatureSnapshot,
			gen.FeatureIntercept,
		},
	}

	if err := entc.Generate(schemaPath,
		genConfig,
		// 扩展自定义模版的生成
		entc.TemplateDir("template"),
	); err != nil {
		log.Fatalf("failed running third ent codegen with hooks and policy: %s", err)
	}
	//// 生成ent 骨架
	//if err := entc.Generate(schemaPath, genConfig, entc.BuildTags("skiphooks", "skippolicy")); err != nil {
	//	log.Fatalf("failed running first ent codegen: %s", err)
	//}
	//// 生成 hooks
	//if err := entc.Generate(schemaPath, genConfig, entc.BuildTags("skippolicy")); err != nil {
	//	log.Fatalf("failed running second ent codegen with hooks: %s", err)
	//}
	// 生成

}
