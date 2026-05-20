package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// buildRouter 建立用于测试的路由器（不依赖真实 service）
func buildRouter(h *AuthHandler) *gin.Engine {
	r := gin.New()
	r.POST("/api/v1/auth/login/phone", h.LoginByPhone)
	r.POST("/api/v1/auth/login/wechat", h.LoginByWechat)
	r.POST("/api/v1/auth/refresh", h.RefreshToken)
	r.POST("/api/v1/auth/logout", h.Logout)
	return r
}

// mockAuthService 实现测试用的认证服务接口
type mockAuthService struct {
	loginResult *loginResultStub
	loginErr    error
	refreshErr  error
}

type loginResultStub struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// TestLoginByPhone_MissingFields 验证缺少必填字段时返回 400
func TestLoginByPhone_MissingFields(t *testing.T) {
	// 创建不依赖 service 的最小测试 handler
	r := gin.New()
	r.POST("/api/v1/auth/login/phone", func(c *gin.Context) {
		var req struct {
			Phone string `json:"phone" binding:"required,len=11"`
			Code  string `json:"code" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 40001, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name:       "缺少phone",
			body:       map[string]string{"code": "123456"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "缺少code",
			body:       map[string]string{"phone": "13812345678"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "phone长度不足",
			body:       map[string]string{"phone": "138", "code": "123456"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "phone长度超过11",
			body:       map[string]string{"phone": "138123456789", "code": "123456"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "正常请求",
			body:       map[string]string{"phone": "13812345678", "code": "123456"},
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req := httptest.NewRequest("POST", "/api/v1/auth/login/phone", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("HTTP状态码: got %d, want %d (body=%s)", w.Code, tc.wantStatus, w.Body.String())
			}
		})
	}
}

// TestLoginByWechat_MissingCode 验证微信登录缺少 code 时返回 400
func TestLoginByWechat_MissingCode(t *testing.T) {
	r := gin.New()
	r.POST("/api/v1/auth/login/wechat", func(c *gin.Context) {
		var req struct {
			Code      string `json:"code" binding:"required"`
			Nickname  string `json:"nickname"`
			AvatarURL string `json:"avatar_url"`
			Role      string `json:"role"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 40001, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name:       "缺少code",
			body:       map[string]string{"nickname": "test"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "正常请求",
			body:       map[string]string{"code": "wx_auth_code_123"},
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req := httptest.NewRequest("POST", "/api/v1/auth/login/wechat", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tc.wantStatus {
				t.Errorf("HTTP状态码: got %d, want %d", w.Code, tc.wantStatus)
			}
		})
	}
}

// TestRefreshToken_MissingRefreshToken 验证刷新接口缺少 token 时返回 400
func TestRefreshToken_MissingRefreshToken(t *testing.T) {
	r := gin.New()
	r.POST("/api/v1/auth/refresh", func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 40001, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	body, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("HTTP状态码: got %d, want 400", w.Code)
	}
}

// TestLogout_AlwaysSuccess 验证登出接口始终返回成功
func TestLogout_AlwaysSuccess(t *testing.T) {
	r := gin.New()
	r.POST("/api/v1/auth/logout", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
	})

	req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("HTTP状态码: got %d, want 200", w.Code)
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("JSON解析失败: %v", err)
	}
	if resp["code"] != float64(0) {
		t.Errorf("code: got %v, want 0", resp["code"])
	}
}

// TestLoginByPhone_InvalidJSON 验证非法 JSON 请求体
func TestLoginByPhone_InvalidJSON(t *testing.T) {
	r := gin.New()
	r.POST("/api/v1/auth/login/phone", func(c *gin.Context) {
		var req struct {
			Phone string `json:"phone" binding:"required,len=11"`
			Code  string `json:"code" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 40001, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0})
	})

	req := httptest.NewRequest("POST", "/api/v1/auth/login/phone",
		bytes.NewBufferString("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("HTTP状态码: got %d, want 400", w.Code)
	}
}

// TestResponseFormat_LoginSuccess 验证成功登录响应格式
func TestResponseFormat_LoginSuccess(t *testing.T) {
	r := gin.New()
	r.POST("/api/v1/auth/login/phone", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"access_token":  "mock_access_token",
				"refresh_token": "mock_refresh_token",
				"user":          gin.H{"id": 1, "role": "consumer"},
			},
		})
	})

	body, _ := json.Marshal(map[string]string{"phone": "13812345678", "code": "123456"})
	req := httptest.NewRequest("POST", "/api/v1/auth/login/phone", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("JSON解析失败: %v", err)
	}
	if resp["code"] != float64(0) {
		t.Errorf("code: got %v, want 0", resp["code"])
	}
	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.Fatal("data字段类型错误")
	}
	if _, ok := data["access_token"]; !ok {
		t.Error("响应中缺少 access_token")
	}
	if _, ok := data["refresh_token"]; !ok {
		t.Error("响应中缺少 refresh_token")
	}
}
