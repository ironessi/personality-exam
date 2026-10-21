package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"PersonalityTest/internal/dao"

	"github.com/gogf/gf/v2/frame/g"
)

// 本文件是 ADR-0007 的 worker：从 Redis 队列取任务，用完整 30 秒超时调 Python，
// 成功后回填 user_test_results.result_data 里的 rag 字段。
//
// worker 的存在理由只有一个：**同步 3 秒装不下一次真实增强**（实测 4.8~9.3 秒）。
// 没有它，那个 3 秒期限就意味着"增强几乎永远拿不到"。

const (
	// brpopTimeoutSeconds 是 BRPOP 的阻塞秒数。设成 5 是为了让 ctx 取消
	// 能在 5 秒内被感知到，同时不至于空转打满 CPU。
	brpopTimeoutSeconds = 5
	// workerRetryDelay 是 Redis 不可用时的退避，别把日志刷爆。
	workerRetryDelay = 3 * time.Second
)

// StartWorker 启动一个后台 goroutine 消费增强队列。
//
// 只在 rag.enable 为真时启动：关掉增强就不该有 worker 连 Redis。
// 注意它**不阻塞**启动流程——队列消费失败不能影响 HTTP 服务起来。
func StartWorker(ctx context.Context, cfg Config) {
	if !cfg.Enable {
		return
	}

	client := NewClient(cfg)
	queue, err := redisClient(cfg)
	if err != nil {
		g.Log().Warningf(ctx, "RAG worker 未启动：连接 Redis 失败: %v", err)
		return
	}

	go func() {
		g.Log().Infof(ctx, "RAG worker 已启动，队列=%s redis=%s/%d",
			cfg.QueueKey, cfg.RedisAddress, cfg.RedisDB)

		for {
			if ctx.Err() != nil {
				g.Log().Info(ctx, "RAG worker 退出")
				return
			}

			vals, err := queue.GroupList().BRPop(ctx, brpopTimeoutSeconds, cfg.QueueKey)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				g.Log().Warningf(ctx, "RAG worker 取任务失败（%.0fs 后重试）: %v",
					workerRetryDelay.Seconds(), err)
				time.Sleep(workerRetryDelay)
				continue
			}
			// BRPop 返回 [key, value]；超时返回空切片。
			if len(vals) < 2 {
				continue
			}

			var job Job
			if err := json.Unmarshal([]byte(vals[1].String()), &job); err != nil {
				// 坏任务只能丢弃：留在队列里会让 worker 每次都崩在同一个上了。
				g.Log().Warningf(ctx, "RAG worker 丢弃无法解析的任务: %v", err)
				continue
			}

			processJob(ctx, client, cfg, job)
		}
	}()
}

// processJob 跑一条任务并回填。
func processJob(ctx context.Context, client *Client, cfg Config, job Job) {
	if job.Request == nil {
		g.Log().Warningf(ctx, "RAG worker 丢弃空请求的任务 record_id=%d", job.RecordID)
		return
	}

	// 这里用 worker 自己的超时（默认 30s），**不是**同步路径的 3 秒期限。
	jobCtx, cancel := context.WithTimeout(ctx, cfg.WorkerTimeout())
	defer cancel()

	start := time.Now()
	resp, err := client.Enhance(jobCtx, job.Request)
	elapsed := time.Since(start).Round(time.Millisecond)

	if err != nil {
		g.Log().Warningf(ctx, "RAG worker: record_id=%d 增强失败（%s）: %v",
			job.RecordID, elapsed, err)
		return
	}
	if !resp.Success {
		g.Log().Warningf(ctx, "RAG worker: record_id=%d 增强返回失败: %s", job.RecordID, resp.Error)
		return
	}

	ragMap := resp.Stored()
	if ragMap == nil {
		g.Log().Warningf(ctx, "RAG worker: record_id=%d 增强结果为空", job.RecordID)
		return
	}

	if err := backfill(ctx, job.RecordID, ragMap); err != nil {
		g.Log().Errorf(ctx, "RAG worker: record_id=%d 回填失败: %v", job.RecordID, err)
		return
	}
	g.Log().Infof(ctx, "RAG worker: record_id=%d 增强完成并回填（%s）", job.RecordID, elapsed)
}

// backfill 把增强结果写回该行的 result_data.rag，并同步顶层 rag_status。
//
// 为什么要读-改-写整份 JSON：result_data 是一整个 JSON 列，MySQL 侧没有
// 稳定的"只改一个子键"的通用写法（JSON_SET 要拼路径，且列类型/版本敏感）。
// ADR-0007 保证了并发安全：同步路径超时后 handler 不再写这一行，
// worker 是唯一的写者。
func backfill(ctx context.Context, recordID int64, ragMap map[string]interface{}) error {
	return updateResultData(ctx, recordID, func(result map[string]interface{}) {
		result["rag"] = ragMap
		// 顶层也镜像一份，前端只看这一个字段就能决定要不要轮询（§6.3）。
		// 值取自 ragMap["status"]，与 Python 的 RagStatus 枚举同源。
		if status, ok := ragMap["status"].(string); ok {
			result["rag_status"] = status
		}
	})
}

// MarkFailed 把顶层 rag_status 标成 failed。
//
// 用于「任务根本没能入队、也没有同步结果」的情况：否则该行会永远停在
// pending，前端会一直轮询一个不会有结果的任务。
func MarkFailed(ctx context.Context, recordID int64) error {
	return updateResultData(ctx, recordID, func(result map[string]interface{}) {
		result["rag_status"] = "failed"
	})
}

func updateResultData(
	ctx context.Context,
	recordID int64,
	mutate func(result map[string]interface{}),
) error {
	one, err := dao.UserTestResults.Ctx(ctx).WherePri(recordID).One()
	if err != nil {
		return fmt.Errorf("读取测评结果失败: %w", err)
	}
	if one.IsEmpty() {
		return fmt.Errorf("测评结果 %d 不存在（可能已被删除）", recordID)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(one["result_data"].String()), &result); err != nil {
		return fmt.Errorf("解析已有 result_data 失败: %w", err)
	}
	mutate(result)

	updated, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("序列化回填后的结果失败: %w", err)
	}

	if _, err := dao.UserTestResults.Ctx(ctx).WherePri(recordID).Update(
		g.Map{"result_data": string(updated)},
	); err != nil {
		return fmt.Errorf("写回 result_data 失败: %w", err)
	}
	return nil
}
