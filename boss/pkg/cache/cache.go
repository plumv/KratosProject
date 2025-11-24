package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)
import "github.com/go-kratos/kratos/v2/log"

// Cache 缓存数据库
type Cache struct {
	client *redis.Client
	log    *log.Helper
}

func NewCache(c *redis.Client, logger log.Logger) *Cache {
	return &Cache{
		client: c,
		log:    log.NewHelper(logger),
	}
}

// Add 把任意对象序列化后存入缓存
func (c *Cache) Add(ctx context.Context, key string, value any, ttl time.Duration) {
	bytes, err := json.Marshal(value)
	if err != nil {
		c.log.Errorf("cache marshal failed, key=%s err=%v", key, err)
		return
	}
	if err := c.client.Set(ctx, key, bytes, ttl).Err(); err != nil {
		c.log.Errorf("cache set failed, key=%s err=%v", key, err)
		return
	}
}

// Delete 删除缓存
func (c *Cache) Delete(ctx context.Context, keys ...string) {
	if len(keys) == 0 {
		return
	}
	if err := c.client.Del(ctx, keys...).Err(); err != nil {
		c.log.Errorf("cache del failed, keys=%v err=%v", keys, err)
		return
	}
	return
}

// Get 读取并反序列化到 dest（必须传指针）
func (c *Cache) Get(ctx context.Context, key string, dest any) error {
	bytes, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		c.log.Errorf("cache get miss, key=%s", key)
		return ErrCacheMiss // 方便上层判断未命中
	}
	if err != nil {
		c.log.Errorf("cache get failed, key=%s err=%v", key, err)
		return err
	}
	if err := json.Unmarshal(bytes, dest); err != nil {
		c.log.Errorf("cache unmarshal failed, key=%s err=%v", key, err)
		return err
	}
	return nil
}

// Update 直接覆盖已有 key（等价于 Add，语义更明确）
func (c *Cache) Update(ctx context.Context, key string, value any, ttl time.Duration) {
	c.Add(ctx, key, value, ttl)
}

// ErrCacheMiss 预定义错误，上层可直接 errors.Is 判断
var ErrCacheMiss = errors.New("cache miss")
