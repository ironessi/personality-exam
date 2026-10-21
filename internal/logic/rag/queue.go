package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/database/gredis"
)

// 本文件是 ADR-0007 的异步兜底：同步 3 秒没等到结果，就把任务丢进 Redis 队列，
// 由 worker 用完整的 30 秒超时慢慢跑，跑完回填 result_data.rag。
//
// 为什么不是「超时就放弃」：实测单次增强 4.8~9.3 秒，3 秒覆盖不了。如果超时即
// 放弃，绝大多数提交都拿不到增强，那这个功能等于没做；而如果放开同步期限，
// 用户就要在提交页等到 9 秒。两条路都不可接受，所以必须两条都留着。

// Job 是队列里的一条任务。
//
// 带上完整 Request 而不是 record_id + 让 worker 去查库：拼请求依赖
// resultMap（含 factors[] 等算出结果），worker 重新查库要再走一遍解析，
// 而且那时读到的可能是已经被改写过的行。请求体本身就是"当时的输入快照"。
type Job struct {
	RecordID  int64           `json:"record_id"`
	TestType  int             `json:"test_type"`
	Request   *EnhanceRequest `json:"request"`
	CreatedAt int64           `json:"created_at"`
}

// redisClient 按 rag.redis.* 构造一个独立客户端。
//
// 不用 g.Redis()：它读的是 redis.default（项目里指向远端 172.16.0.203 的
// 业务缓存），把增强队列混进业务缓存里会互相干扰，也让 flush 变得危险。
func redisClient(cfg Config) (*gredis.Redis, error) {
	return gredis.New(&gredis.Config{
		Address: cfg.RedisAddress,
		Db:      cfg.RedisDB,
		Pass:    cfg.RedisPass,
	})
}

// Enqueue 把任务推进队列（LPUSH + worker 用 BRPOP，天然 FIFO）。
//
// 返回错误只是为了让调用方记日志；提交接口**不得**因此失败。
func (p *Pending) Enqueue(ctx context.Context, recordID int64) error {
	job := Job{
		RecordID:  recordID,
		TestType:  p.TestType,
		Request:   p.Request,
		CreatedAt: time.Now().Unix(),
	}
	raw, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("序列化增强任务失败: %w", err)
	}

	// 队列本身也要有超时：Redis 挂了不该把提交接口拖住。
	pushCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	rdb, err := redisClient(p.cfg)
	if err != nil {
		return fmt.Errorf("连接增强队列 Redis 失败: %w", err)
	}
	if _, err := rdb.GroupList().LPush(pushCtx, p.cfg.QueueKey, raw); err != nil {
		return fmt.Errorf("增强任务入队失败: %w", err)
	}
	return nil
}
