package rag

import (
	"context"
	"testing"
	"time"

	// 测试进程不会跑 main.go，必须自己注册 redis adapter，
	// 否则 gredis.New 报 "redis adapter is not set"。
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
)

func fpReq(testName string, strength float64) *EnhanceRequest {
	s := strength
	return &EnhanceRequest{
		SchemaVersion: 1,
		Test:          TestInfo{Code: 13, Key: "mbti", Name: testName, Domains: []string{"typology", "common"}, Policy: "full"},
		Measures: []Measure{{
			Key: "EI", Label: "外倾 / 内倾", Kind: "dimension", Score: 15,
			Threshold: fptr(12), Strength: &s, Pole: "E",
		}},
	}
}

func fptr(v float64) *float64 { return &v }

// TestCacheFingerprintIgnoresCosmeticFields 盯的是「缓存该命中的时候要命中」：
// 测评名、记录 ID、用户 ID 都不影响解读内容，混进指纹只会让命中率归零。
func TestCacheFingerprintIgnoresCosmeticFields(t *testing.T) {
	a, err := cacheFingerprint(fpReq("16型人格测试（MBTI）", 0.5))
	if err != nil {
		t.Fatal(err)
	}
	b, err := cacheFingerprint(fpReq("MBTI", 0.5))
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Error("只改测评名不应改变指纹")
	}
}

// TestCacheFingerprintSensitiveToScores 反过来：档案变了必须换键，
// 否则会把别人的解读返回给这个人。
func TestCacheFingerprintSensitiveToScores(t *testing.T) {
	a, _ := cacheFingerprint(fpReq("MBTI", 0.5))
	b, _ := cacheFingerprint(fpReq("MBTI", 0.51))
	if a == b {
		t.Error("strength 变了，指纹必须变")
	}

	c := fpReq("MBTI", 0.5)
	c.Measures[0].Score = 16
	d, _ := cacheFingerprint(c)
	if a == d {
		t.Error("score 变了，指纹必须变")
	}
}

func TestCacheFingerprintEmptyRequest(t *testing.T) {
	key, err := cacheFingerprint(nil)
	if err != nil || key != "" {
		t.Errorf("nil 请求应返回空键，实得 (%q, %v)", key, err)
	}
}

// TestCacheRoundTrip 需要一个可用的 Redis；不可用时跳过而不是失败——
// 缓存是可选加速，不该让单测依赖本机环境。
func TestCacheRoundTrip(t *testing.T) {
	cfg := Config{
		RedisAddress:    "127.0.0.1:6379",
		RedisDB:         5,
		QueueKey:        "rag:queue:test",
		CacheTTLSeconds: 60,
	}
	ctx := context.Background()

	rdb, err := redisClient(cfg)
	if err != nil {
		t.Skipf("Redis 不可用，跳过：%v", err)
	}
	pingCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	if _, err := rdb.Do(pingCtx, "PING"); err != nil {
		t.Skipf("Redis 不可用，跳过：%v", err)
	}
	_ = rdb.Close(ctx)

	p := &Pending{TestType: 13, Request: fpReq("16型人格测试（MBTI）", 0.5), cfg: cfg}
	stored := map[string]interface{}{"status": "ready", "enhanced_text": "缓存文本"}

	p.SaveCache(ctx, stored)
	got := p.TryCache(ctx)
	if got == nil {
		t.Fatal("缓存未命中")
	}
	if got["enhanced_text"] != "缓存文本" {
		t.Errorf("缓存内容不对：%v", got)
	}

	// 清理，避免污染本机 redis db5
	if rdb2, err := redisClient(cfg); err == nil {
		if key, err := cacheFingerprint(p.Request); err == nil {
			_, _ = rdb2.Del(ctx, cacheKeyPrefix+key)
		}
		_ = rdb2.Close(ctx)
	}
}
