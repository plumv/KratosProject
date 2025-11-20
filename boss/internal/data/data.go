package data

import (
	"boss/internal/biz"
	"boss/internal/conf"
	"boss/pkg/ent"
	"boss/pkg/ent/hook"
	"boss/pkg/ent/intercept"
	"boss/pkg/ent/migrate"
	"context"
	"fmt"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewTransaction, NewUserRepo, NewRoleRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	db  *ent.Database
	rdb *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	lh := log.NewHelper(logger)
	// 构建Driver
	drv, err := sql.Open(c.Database.Driver, c.Database.Source)
	// 添加日志
	sqlDrv := dialect.DebugWithContext(drv, func(ctx context.Context, i ...interface{}) {
		lh.WithContext(ctx).Info(i...)
		tracer := otel.Tracer("ent.")
		kind := trace.SpanKindServer
		_, span := tracer.Start(ctx,
			"Query",
			trace.WithAttributes(
				attribute.String("sql", fmt.Sprint(i...)),
			),
			trace.WithSpanKind(kind),
		)
		span.End()
	})
	// 构建客户端
	db := ent.NewDatabase(ent.Driver(sqlDrv))
	if err != nil {
		log.Errorf("failed opening connection to sqlite: %v", err)
		return nil, nil, err
	}
	// 添加软删除逻辑
	db.Client.Intercept(intercept.SoftDeleteInterceptor)
	db.Client.Use(hook.SoftDeleteHook)

	// 创建表
	if err := db.Client.Schema.Create(
		context.Background(),
		migrate.WithGlobalUniqueID(true),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		log.Errorf("failed creating schema resources: %v", err)
		return nil, nil, err
	}
	// 构建redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DB:           int(c.Redis.Db),
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
	})
	rdb.AddHook(redisotel.TracingHook{})

	// 构建
	d := &Data{
		db:  db,
		rdb: rdb,
	}
	return d, func() {
		lh.Info("closing the data resources")
		if err := d.db.Close(); err != nil {
			log.Error(err)
		}
	}, nil
}

func NewTransaction(data *Data) biz.Transaction {
	return data.db
}
