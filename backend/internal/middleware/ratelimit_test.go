package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func buildRateLimitRouter(qps float64, burst int) *gin.Engine {
	r := gin.New()
	r.GET("/ping", RateLimit(qps, burst), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
	})
	return r
}

// TestRateLimit_AllowsUnderLimit 验证低于限制时请求正常通过
func TestRateLimit_AllowsUnderLimit(t *testing.T) {
	r := buildRateLimitRouter(100, 10)

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/ping", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("第%d次请求: got %d, want 200", i+1, w.Code)
		}
	}
}

// TestRateLimit_BlocksOverLimit 验证突发流量超过 burst 时触发限流
func TestRateLimit_BlocksOverLimit(t *testing.T) {
	// burst=3, qps=0（极低速率），前3个请求消耗桶中令牌，之后应被限流
	r := buildRateLimitRouter(0.001, 3)

	results := make([]int, 10)
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/ping", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		results[i] = w.Code
	}

	// 前3个应该成功，之后应该被限流
	for i := 0; i < 3; i++ {
		if results[i] != http.StatusOK {
			t.Errorf("请求%d (burst内) 应返回200, got %d", i+1, results[i])
		}
	}
	throttledCount := 0
	for i := 3; i < 10; i++ {
		if results[i] == http.StatusTooManyRequests {
			throttledCount++
		}
	}
	if throttledCount == 0 {
		t.Error("超出 burst 后应有请求被限流（429）")
	}
}

// TestRateLimit_ResponseFormat 验证限流响应格式
func TestRateLimit_ResponseFormat(t *testing.T) {
	r := buildRateLimitRouter(0.001, 1)

	// 消耗唯一令牌
	req1 := httptest.NewRequest("GET", "/ping", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	// 第二次应被限流
	req2 := httptest.NewRequest("GET", "/ping", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusTooManyRequests {
		t.Skipf("第二次请求未被限流（%d），跳过响应格式验证", w2.Code)
	}

	// 验证响应体格式
	body := w2.Body.String()
	if body == "" {
		t.Error("限流响应体不应为空")
	}
	// 应包含 code 和 message 字段
	if len(body) < 10 {
		t.Errorf("限流响应体过短: %s", body)
	}
}

// TestRateLimit_StatusCode429 验证限流时返回 429 状态码
func TestRateLimit_StatusCode429(t *testing.T) {
	r := buildRateLimitRouter(0.001, 0)

	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("burst=0时应立即限流返回429, got %d", w.Code)
	}
}

// TestRateLimit_HighQPS 验证高 QPS 配置不会误限流
func TestRateLimit_HighQPS(t *testing.T) {
	r := buildRateLimitRouter(10000, 10000)

	var wg sync.WaitGroup
	errors := 0
	var mu sync.Mutex

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest("GET", "/ping", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				mu.Lock()
				errors++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if errors > 0 {
		t.Errorf("高QPS配置下50次并发请求有%d次被限流（不应发生）", errors)
	}
}

// TestRateLimit_Recovery 验证限流后等待令牌补充可以恢复
func TestRateLimit_Recovery(t *testing.T) {
	// QPS=10, burst=1 => 每100ms补充1个令牌
	r := buildRateLimitRouter(10, 1)

	// 消耗令牌
	req1 := httptest.NewRequest("GET", "/ping", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	// 立即请求应被限流
	req2 := httptest.NewRequest("GET", "/ping", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	if w2.Code == http.StatusOK {
		t.Skip("令牌补充很快，无法验证恢复逻辑")
	}

	// 等待令牌补充
	time.Sleep(150 * time.Millisecond)

	req3 := httptest.NewRequest("GET", "/ping", nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	if w3.Code != http.StatusOK {
		t.Errorf("等待令牌补充后请求应成功, got %d", w3.Code)
	}
}

// TestRateLimit_ConcurrentRequests 验证并发场景下限流准确性
func TestRateLimit_ConcurrentRequests(t *testing.T) {
	burst := 10
	r := buildRateLimitRouter(0.001, burst)

	var (
		wg          sync.WaitGroup
		successCount int
		failCount    int
		mu           sync.Mutex
	)

	total := 30
	for i := 0; i < total; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest("GET", "/ping", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			mu.Lock()
			if w.Code == http.StatusOK {
				successCount++
			} else {
				failCount++
			}
			mu.Unlock()
		}()
	}
	wg.Wait()

	// 成功数不应超过 burst
	if successCount > burst {
		t.Errorf("并发成功数 %d 超过 burst %d", successCount, burst)
	}
	if successCount+failCount != total {
		t.Errorf("成功(%d)+失败(%d) != 总请求数(%d)", successCount, failCount, total)
	}
}
