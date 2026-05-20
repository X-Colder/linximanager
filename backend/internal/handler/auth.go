package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/linximanager/backend/internal/pkg/errcode"
	"github.com/linximanager/backend/internal/pkg/response"
	"github.com/linximanager/backend/internal/pkg/validator"
	"github.com/linximanager/backend/internal/service"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

type phoneLoginReq struct {
	Phone string `json:"phone" binding:"required,len=11"`
	Code  string `json:"code" binding:"required"`
}

func (h *AuthHandler) LoginByPhone(c *gin.Context) {
	var req phoneLoginReq
	if !validator.BindJSON(c, &req) {
		return
	}
	result, err := h.authSvc.LoginByPhone(c.Request.Context(), req.Phone, req.Code)
	if err != nil {
		if err.Error() == "account frozen" {
			response.Fail(c, errcode.ErrMerchantFrozen)
			return
		}
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OK(c, result)
}

type wechatLoginReq struct {
	Code      string `json:"code" binding:"required"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	Role      string `json:"role"` // merchant/consumer
}

func (h *AuthHandler) LoginByWechat(c *gin.Context) {
	var req wechatLoginReq
	if !validator.BindJSON(c, &req) {
		return
	}
	// 生产中此处用 code 换取 openid/unionid
	openid := "mock_openid_" + req.Code
	role := req.Role
	if role == "" {
		role = "consumer"
	}
	result, err := h.authSvc.LoginByWechat(c.Request.Context(), openid, "", req.Nickname, req.AvatarURL, role)
	if err != nil {
		response.Fail(c, errcode.ErrInternal)
		return
	}
	response.OK(c, result)
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req refreshReq
	if !validator.BindJSON(c, &req) {
		return
	}
	result, err := h.authSvc.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Fail(c, errcode.ErrTokenInvalid)
		return
	}
	response.OK(c, result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// JWT 无状态，客户端丢弃 token 即可；生产中可加黑名单
	response.OK(c, nil)
}
