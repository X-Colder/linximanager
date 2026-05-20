package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linximanager/backend/internal/pkg/errcode"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type PageData struct {
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	List     any   `json:"list"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:    errcode.Success,
		Message: "success",
		Data:    data,
	})
}

func OKPage(c *gin.Context, total int64, page, pageSize int, list any) {
	c.JSON(http.StatusOK, Response{
		Code:    errcode.Success,
		Message: "success",
		Data: PageData{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
			List:     list,
		},
	})
}

func Fail(c *gin.Context, code int) {
	httpStatus := codeToHTTPStatus(code)
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: errcode.Message(code),
	})
}

func FailMsg(c *gin.Context, code int, msg string) {
	httpStatus := codeToHTTPStatus(code)
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: msg,
	})
}

func codeToHTTPStatus(code int) int {
	switch {
	case code == errcode.ErrUnauthorized || code == errcode.ErrTokenExpired || code == errcode.ErrTokenInvalid:
		return http.StatusUnauthorized
	case code == errcode.ErrForbidden || code == errcode.ErrMerchantFrozen || code == errcode.ErrMerchantAuditPend:
		return http.StatusForbidden
	case code == errcode.ErrNotFound || code == errcode.ErrUserNotFound || code == errcode.ErrMerchantNotFound ||
		code == errcode.ErrProductNotFound || code == errcode.ErrOrderNotFound:
		return http.StatusNotFound
	case code == errcode.ErrTooManyRequests:
		return http.StatusTooManyRequests
	case code >= 50000:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}
