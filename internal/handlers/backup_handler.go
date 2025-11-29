package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"itab-backend/internal/database"
	"itab-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// ListBackups 获取备份列表
func ListBackups(c *gin.Context) {
	userID := c.GetUint("user_id")
	isAdmin := c.GetBool("is_admin")

	var backups []models.Backup
	query := database.DB.Preload("User").Select("id, name, size, sync_count, user_id, created_at, updated_at")

	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&backups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取备份列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": backups})
}

// GetBackup 获取备份详情
func GetBackup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的备份ID"})
		return
	}

	userID := c.GetUint("user_id")
	isAdmin := c.GetBool("is_admin")

	var backup models.Backup
	if err := database.DB.Preload("User").First(&backup, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "备份不存在"})
		return
	}

	// 非管理员只能查看自己的备份
	if !isAdmin && backup.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权查看此备份"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": backup})
}

// DeleteBackup 删除备份
func DeleteBackup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的备份ID"})
		return
	}

	userID := c.GetUint("user_id")
	isAdmin := c.GetBool("is_admin")

	var backup models.Backup
	if err := database.DB.First(&backup, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "备份不存在"})
		return
	}

	// 非管理员只能删除自己的备份
	if !isAdmin && backup.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除此备份"})
		return
	}

	if err := database.DB.Delete(&backup).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除备份失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "备份删除成功"})
}

// DownloadBackup 下载备份数据
func DownloadBackup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的备份ID"})
		return
	}

	userID := c.GetUint("user_id")
	isAdmin := c.GetBool("is_admin")

	var backup models.Backup
	if err := database.DB.First(&backup, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "备份不存在"})
		return
	}

	// 非管理员只能下载自己的备份
	if !isAdmin && backup.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权下载此备份"})
		return
	}

	// 解析data为对象
	var backupData interface{}
	if err := json.Unmarshal([]byte(backup.Data), &backupData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析备份数据失败"})
		return
	}

	// 返回完整的导出格式
	c.JSON(http.StatusOK, gin.H{
		"version":            "2.0",
		"exportDate":         backup.UpdatedAt,
		"passwordsEncrypted": backup.PasswordsEncrypted,
		"data":               backupData,
	})
}
