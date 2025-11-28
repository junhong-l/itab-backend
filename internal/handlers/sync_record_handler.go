package handlers

import (
	"net/http"
	"strconv"
	"time"

	"itab-backend/internal/database"
	"itab-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// ListSyncRecords 获取同步记录列表
func ListSyncRecords(c *gin.Context) {
	userID := c.GetUint("user_id")
	isAdmin := c.GetBool("is_admin")

	var records []models.SyncRecord
	query := database.DB.Preload("User").Order("created_at DESC")

	if !isAdmin {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取同步记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": records})
}

// CleanSyncRecords 清理同步记录
type CleanRequest struct {
	Days int `json:"days"` // 0表示清理全部
}

func CleanSyncRecords(c *gin.Context) {
	var req CleanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	query := database.DB.Where("1 = 1")

	if req.Days > 0 {
		cutoff := time.Now().AddDate(0, 0, -req.Days)
		query = database.DB.Where("created_at < ?", cutoff)
	}

	result := query.Delete(&models.SyncRecord{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "清理记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "清理成功",
		"deleted": result.RowsAffected,
	})
}

// GetSyncStats 获取同步统计
func GetSyncStats(c *gin.Context) {
	var totalRecords int64
	var uploadCount int64
	var downloadCount int64

	database.DB.Model(&models.SyncRecord{}).Count(&totalRecords)
	database.DB.Model(&models.SyncRecord{}).Where("trans_type = ?", "upload").Count(&uploadCount)
	database.DB.Model(&models.SyncRecord{}).Where("trans_type = ?", "download").Count(&downloadCount)

	c.JSON(http.StatusOK, gin.H{
		"total":     totalRecords,
		"uploads":   uploadCount,
		"downloads": downloadCount,
	})
}

// DeleteSyncRecord 删除单条同步记录
func DeleteSyncRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	if err := database.DB.Delete(&models.SyncRecord{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
