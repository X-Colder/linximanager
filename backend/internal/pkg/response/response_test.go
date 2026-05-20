package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/linximanager/backend/internal/pkg/errcode"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestContext(method, path string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(method, path, nil)
	c.Request = req
	return c, w
}

func TestOK_WithData(t *testing.T) {
	c, w := newTestContext("GET", "/test")
	OK(c, map[string]string{"key": "value"})

	if w.Code != http.StatusOK {
		t.Errorf("HTTP状态码: got %d, want %d", w.Code, http.StatusOK)
	}
	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("JSON解析失败: %v", err)
	}
	if resp.Code != errcode.Success {
		t.Errorf("code: got %d, want %d", resp.Code, errcode.Success)
	}
	if resp.Message != "success" {
		t.Errorf("message: got %s, want success", resp.Message)
	}
	if resp.Data == nil {
		t.Error("data不应为nil")
	}
}

func TestOK_WithNilData(t *testing.T) {
	c, w := newTestContext("GET", "/test")
	OK(c, nil)

	if w.Code != http.StatusOK {
		t.Errorf("HTTP状态码: got %d, want %d", w.Code, http.StatusOK)
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("JSON解析失败: %v", err)
	}
	if _, ok := resp["data"]; ok {
		t.Error("data字段应被omitempty省略")
	}
}

func TestOKPage(t *testing.T) {
	c, w := newTestContext("GET", "/test")
	items := []string{"a", "b", "c"}
	OKPage(c, 100, 2, 20, items)

	if w.Code != http.StatusOK {
		t.Errorf("HTTP状态码: got %d, want %d", w.Code, http.StatusOK)
	}
	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("JSON解析失败: %v", err)
	}
	if resp.Code != errcode.Success {
		t.Errorf("code: got %d, want %d", resp.Code, errcode.Success)
	}

	dataBytes, _ := json.Marshal(resp.Data)
	var pageData PageData
	if err := json.Unmarshal(dataBytes, &pageData); err != nil {
		t.Fatalf("PageData解析失败: %v", err)
	}
	if pageData.Total != 100 {
		t.Errorf("Total: got %d, want 100", pageData.Total)
	}
	if pageData.Page != 2 {
		t.Errorf("Page: got %d, want 2", pageData.Page)
	}
	if pageData.PageSize != 20 {
		t.Errorf("PageSize: got %d, want 20", pageData.PageSize)
	}
}

func TestFail_HTTPStatusMapping(t *testing.T) {
	tests := []struct {
		name           string
		code           int
		wantHTTPStatus int
	}{
		{"未认证", errcode.ErrUnauthorized, http.StatusUnauthorized},
		{"Token过期", errcode.ErrTokenExpired, http.StatusUnauthorized},
		{"Token无效", errcode.ErrTokenInvalid, http.StatusUnauthorized},
		{"无权限", errcode.ErrForbidden, http.StatusForbidden},
		{"商家冻结", errcode.ErrMerchantFrozen, http.StatusForbidden},
		{"商家待审核", errcode.ErrMerchantAuditPend, http.StatusForbidden},
		{"资源不存在", errcode.ErrNotFound, http.StatusNotFound},
		{"用户不存在", errcode.ErrUserNotFound, http.StatusNotFound},
		{"商家不存在", errcode.ErrMerchantNotFound, http.StatusNotFound},
		{"商品不存在", errcode.ErrProductNotFound, http.StatusNotFound},
		{"订单不存在", errcode.ErrOrderNotFound, http.StatusNotFound},
		{"请求过频", errcode.ErrTooManyRequests, http.StatusTooManyRequests},
		{"内部错误", errcode.ErrInternal, http.StatusInternalServerError},
		{"数据库错误", errcode.ErrDB, http.StatusInternalServerError},
		{"参数无效", errcode.ErrParamInvalid, http.StatusBadRequest},
		{"库存不足", errcode.ErrStockInsufficient, http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c, w := newTestContext("GET", "/test")
			Fail(c, tc.code)
			if w.Code != tc.wantHTTPStatus {
				t.Errorf("errcode %d -> HTTP %d, want %d", tc.code, w.Code, tc.wantHTTPStatus)
			}
			var resp Response
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("JSON解析失败: %v", err)
			}
			if resp.Code != tc.code {
				t.Errorf("业务码: got %d, want %d", resp.Code, tc.code)
			}
			if resp.Message == "" {
				t.Error("message不应为空")
			}
		})
	}
}

func TestFailMsg(t *testing.T) {
	c, w := newTestContext("GET", "/test")
	FailMsg(c, errcode.ErrParamInvalid, "自定义错误消息")

	if w.Code != http.StatusBadRequest {
		t.Errorf("HTTP状态码: got %d, want %d", w.Code, http.StatusBadRequest)
	}
	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("JSON解析失败: %v", err)
	}
	if resp.Message != "自定义错误消息" {
		t.Errorf("message: got %s, want 自定义错误消息", resp.Message)
	}
}
