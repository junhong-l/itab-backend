package middleware

import (
	"net/http"
	"strings"

	"itab-backend/internal/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息"})
			c.Abort()
			return
		}

		// 解析Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证格式错误"})
			c.Abort()
			return
		}

		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("is_admin", claims.IsAdmin)
		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// AccessKeyMiddleware 访问密钥认证中间件
func AccessKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessKey := c.GetHeader("x-access-key")
		secretKey := c.GetHeader("x-secret-key")

		if accessKey == "" || secretKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供访问密钥"})
			c.Abort()
			return
		}

		ak, err := auth.ValidateAccessKey(accessKey, secretKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", ak.UserID)
		c.Set("username", ak.User.Username)
		c.Set("access_key_id", ak.ID)
		c.Set("access_key", ak.AccessKey)
		c.Next()
	}
}
