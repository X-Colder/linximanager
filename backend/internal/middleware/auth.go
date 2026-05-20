package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/linximanager/backend/internal/pkg/jwt"
	"github.com/linximanager/backend/internal/pkg/errcode"
	"github.com/linximanager/backend/internal/pkg/response"
)

const (
	CtxUID  = "uid"
	CtxRole = "role"
	CtxMID  = "mid"
)

// Auth 验证JWT，解析并注入上下文
func Auth(accessSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Fail(c, errcode.ErrUnauthorized)
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Fail(c, errcode.ErrUnauthorized)
			c.Abort()
			return
		}

		claims, err := jwtpkg.ParseAccessToken(parts[1], accessSecret)
		if err != nil {
			if err == jwtpkg.ErrTokenExpired {
				response.Fail(c, errcode.ErrTokenExpired)
			} else {
				response.Fail(c, errcode.ErrTokenInvalid)
			}
			c.Abort()
			return
		}

		c.Set(CtxUID, claims.UID)
		c.Set(CtxRole, claims.Role)
		c.Set(CtxMID, claims.MID)
		c.Next()
	}
}

// RequireRole 检查角色权限，允许多个角色
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		role, _ := c.Get(CtxRole)
		roleStr, _ := role.(string)
		if _, ok := allowed[roleStr]; !ok {
			response.Fail(c, errcode.ErrForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetUID 从上下文中取当前用户ID
func GetUID(c *gin.Context) int64 {
	v, _ := c.Get(CtxUID)
	uid, _ := v.(int64)
	return uid
}

// GetMID 从上下文中取商家ID
func GetMID(c *gin.Context) int64 {
	v, _ := c.Get(CtxMID)
	mid, _ := v.(int64)
	return mid
}

// GetRole 从上下文中取角色
func GetRole(c *gin.Context) string {
	v, _ := c.Get(CtxRole)
	role, _ := v.(string)
	return role
}
