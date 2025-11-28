package handlers

import (
	"net/http"
	"strconv"
	"time"

	"itab-backend/internal/auth"
	"itab-backend/internal/database"
	"itab-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// CreateKeyRequest 创建密钥请求
type CreateKeyRequest struct {
	ExpireDays int `json:"expire_days"` // 0表示永久
}

// ListKeys 获取密钥列表
func ListKeys(c *gin.Context) {
	userID := c.GetUint("user_id")
	isAdmin := c.GetBool("is_admin")

	var keys []models.AccessKey
	query := database.DB.Preload("User")

	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&keys).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取密钥列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": keys})
}

// CreateKey 创建密钥
func CreateKey(c *gin.Context) {
	var req CreateKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.ExpireDays = 30 // 默认30天
	}

	userID := c.GetUint("user_id")
	accessKey, secretKey := auth.GenerateAccessKey()

	key := &models.AccessKey{
		AccessKey: accessKey,
		SecretKey: secretKey,
		UserID:    userID,
	}

	// 设置过期时间
	if req.ExpireDays > 0 {
		expiresAt := time.Now().AddDate(0, 0, req.ExpireDays)
		key.ExpiresAt = &expiresAt
	}

	if err := database.DB.Create(key).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建密钥失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "密钥创建成功",
		"data":    key,
	})
}

// DeleteKey 删除密钥
func DeleteKey(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的密钥ID"})
		return
	}

	userID := c.GetUint("user_id")
	isAdmin := c.GetBool("is_admin")

	var key models.AccessKey
	if err := database.DB.First(&key, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "密钥不存在"})
		return
	}

	// 非管理员只能删除自己的密钥
	if !isAdmin && key.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除此密钥"})
		return
	}

	if err := database.DB.Delete(&key).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除密钥失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密钥删除成功"})
}

// ExpireKey 使密钥过期
func ExpireKey(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的密钥ID"})
		return
	}

	userID := c.GetUint("user_id")
	isAdmin := c.GetBool("is_admin")

	var key models.AccessKey
	if err := database.DB.First(&key, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "密钥不存在"})
		return
	}

	// 非管理员只能操作自己的密钥
	if !isAdmin && key.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作此密钥"})
		return
	}

	key.IsExpired = true
	if err := database.DB.Save(&key).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "操作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密钥已过期"})
}
