package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/linximanager/backend/internal/pkg/jwt"
)

func init() {
	gin.SetMode(gin.TestMode)
}

const testSecret = "test-access-secret-for-middleware"

func buildAuthRouter(secret string) *gin.Engine {
	r := gin.New()
	r.GET("/protected", Auth(secret), func(c *gin.Context) {
		uid := GetUID(c)
		role := GetRole(c)
		mid := GetMID(c)
		c.JSON(http.StatusOK, gin.H{"uid": uid, "role": role, "mid": mid})
	})
	return r
}

func buildRoleRouter(secret string, roles ...string) *gin.Engine {
	r := gin.New()
	r.GET("/admin-only", Auth(secret), RequireRole(roles...), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func makeToken(uid int64, role string, mid int64, expire time.Duration) string {
	token, _ := jwtpkg.GenerateAccessToken(uid, role, mid, testSecret, expire)
	return token
}

// TestAuth_NoToken 验证无 Token 时返回 401
func TestAuth_NoToken(t *testing.T) {
	r := buildAuthRouter(testSecret)
	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP状态码: got %d, want 401", w.Code)
	}
}

// TestAuth_InvalidBearerFormat 验证 Authorization 格式错误时返回 401
func TestAuth_InvalidBearerFormat(t *testing.T) {
	tests := []struct {
		name   string
		header string
	}{
		{"没有Bearer前缀", "just-a-token"},
		{"格式缺少空格", "Bearertoken"},
		{"空字符串", ""},
	}

	r := buildAuthRouter(testSecret)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/protected", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != http.StatusUnauthorized {
				t.Errorf("HTTP状态码: got %d, want 401", w.Code)
			}
		})
	}
}

// TestAuth_ExpiredToken 验证过期 Token 返回 401 且错误码为 40101
func TestAuth_ExpiredToken(t *testing.T) {
	token := makeToken(1, "consumer", 0, -time.Second)
	r := buildAuthRouter(testSecret)
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP状态码: got %d, want 401", w.Code)
	}
	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"] != float64(40101) {
		t.Errorf("业务码: got %v, want 40101", resp["code"])
	}
}

// TestAuth_InvalidToken 验证无效 Token 返回 401 且错误码为 40102
func TestAuth_InvalidToken(t *testing.T) {
	r := buildAuthRouter(testSecret)
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP状态码: got %d, want 401", w.Code)
	}
	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"] != float64(40102) {
		t.Errorf("业务码: got %v, want 40102", resp["code"])
	}
}

// TestAuth_ValidToken 验证有效 Token 可以通过中间件并注入上下文
func TestAuth_ValidToken(t *testing.T) {
	token := makeToken(42, "merchant", 100, time.Hour)
	r := buildAuthRouter(testSecret)
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("HTTP状态码: got %d, want 200 (body=%s)", w.Code, w.Body.String())
	}
	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["uid"] != float64(42) {
		t.Errorf("uid: got %v, want 42", resp["uid"])
	}
	if resp["role"] != "merchant" {
		t.Errorf("role: got %v, want merchant", resp["role"])
	}
	if resp["mid"] != float64(100) {
		t.Errorf("mid: got %v, want 100", resp["mid"])
	}
}

// TestAuth_WrongSecret 验证使用错误 secret 签名的 Token 被拒绝
func TestAuth_WrongSecret(t *testing.T) {
	token, _ := jwtpkg.GenerateAccessToken(1, "consumer", 0, "other-secret", time.Hour)
	r := buildAuthRouter(testSecret)
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP状态码: got %d, want 401", w.Code)
	}
}

// TestRequireRole_Allowed 验证角色匹配时放行
func TestRequireRole_Allowed(t *testing.T) {
	token := makeToken(1, "admin", 0, time.Hour)
	r := buildRoleRouter(testSecret, "admin", "super_admin")
	req := httptest.NewRequest("GET", "/admin-only", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("HTTP状态码: got %d, want 200", w.Code)
	}
}

// TestRequireRole_Forbidden 验证角色不匹配时返回 403
func TestRequireRole_Forbidden(t *testing.T) {
	token := makeToken(1, "consumer", 0, time.Hour)
	r := buildRoleRouter(testSecret, "admin")
	req := httptest.NewRequest("GET", "/admin-only", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("HTTP状态码: got %d, want 403", w.Code)
	}
	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"] != float64(40300) {
		t.Errorf("业务码: got %v, want 40300", resp["code"])
	}
}

// TestRequireRole_MultipleRoles 验证多角色允许
func TestRequireRole_MultipleRoles(t *testing.T) {
	roles := []string{"merchant", "staff"}
	tests := []struct {
		role       string
		wantStatus int
	}{
		{"merchant", http.StatusOK},
		{"staff", http.StatusOK},
		{"consumer", http.StatusForbidden},
		{"admin", http.StatusForbidden},
	}

	for _, tc := range tests {
		t.Run(tc.role, func(t *testing.T) {
			token := makeToken(1, tc.role, 0, time.Hour)
			r := buildRoleRouter(testSecret, roles...)
			req := httptest.NewRequest("GET", "/admin-only", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tc.wantStatus {
				t.Errorf("role=%s: HTTP状态码 got %d, want %d", tc.role, w.Code, tc.wantStatus)
			}
		})
	}
}

// TestGetUID_GetRole_GetMID 验证上下文工具函数
func TestContextHelpers(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 未设置时应返回零值
	if uid := GetUID(c); uid != 0 {
		t.Errorf("未设置UID: got %d, want 0", uid)
	}
	if role := GetRole(c); role != "" {
		t.Errorf("未设置Role: got %s, want empty", role)
	}
	if mid := GetMID(c); mid != 0 {
		t.Errorf("未设置MID: got %d, want 0", mid)
	}

	// 设置后验证
	c.Set(CtxUID, int64(99))
	c.Set(CtxRole, "admin")
	c.Set(CtxMID, int64(200))

	if uid := GetUID(c); uid != 99 {
		t.Errorf("UID: got %d, want 99", uid)
	}
	if role := GetRole(c); role != "admin" {
		t.Errorf("Role: got %s, want admin", role)
	}
	if mid := GetMID(c); mid != 200 {
		t.Errorf("MID: got %d, want 200", mid)
	}
}

// TestAuth_AbortNotCallNext 验证认证失败后不继续执行后续 handler
func TestAuth_AbortNotCallNext(t *testing.T) {
	handlerCalled := false
	r := gin.New()
	r.GET("/test", Auth(testSecret), func(c *gin.Context) {
		handlerCalled = true
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	// 不设置 Authorization
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if handlerCalled {
		t.Error("认证失败后不应执行业务 handler")
	}
	if w.Code != http.StatusUnauthorized {
		t.Errorf("HTTP状态码: got %d, want 401", w.Code)
	}
}

// Suppress fmt import warning
var _ = fmt.Sprintf
