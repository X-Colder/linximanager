package validator

import (
	"github.com/gin-gonic/gin"
	"github.com/linximanager/backend/internal/pkg/errcode"
	"github.com/linximanager/backend/internal/pkg/response"
)

// BindJSON 绑定JSON请求体，失败时写响应并返回false
func BindJSON(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		response.FailMsg(c, errcode.ErrParamInvalid, err.Error())
		return false
	}
	return true
}

// BindQuery 绑定Query参数，失败时写响应并返回false
func BindQuery(c *gin.Context, obj any) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		response.FailMsg(c, errcode.ErrParamInvalid, err.Error())
		return false
	}
	return true
}

// BindUri 绑定URI参数，失败时写响应并返回false
func BindUri(c *gin.Context, obj any) bool {
	if err := c.ShouldBindUri(obj); err != nil {
		response.FailMsg(c, errcode.ErrParamInvalid, err.Error())
		return false
	}
	return true
}

// PageParam 通用分页参数
type PageParam struct {
	Page     int `form:"page" binding:"min=1"`
	PageSize int `form:"page_size" binding:"min=1,max=100"`
}

func (p *PageParam) SetDefault() {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.PageSize == 0 {
		p.PageSize = 20
	}
}

func (p *PageParam) Offset() int {
	return (p.Page - 1) * p.PageSize
}
