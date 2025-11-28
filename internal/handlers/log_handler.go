package handlers

import (
	"net/http"

	"itab-backend/internal/logger"

	"github.com/gin-gonic/gin"
)

// GetLogFiles 获取日志文件列表
func GetLogFiles(c *gin.Context) {
	files, err := logger.GetLogFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": files})
}

// CleanLogsRequest 清理日志请求
type CleanLogsRequest struct {
	Days int `json:"days"` // 清理多少天前的日志，0表示清理所有（除今天）
}

// CleanLogs 清理日志文件
func CleanLogs(c *gin.Context) {
	var req CleanLogsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	deleted, err := logger.CleanLogs(req.Days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "清理成功",
		"deleted": deleted,
	})
}
