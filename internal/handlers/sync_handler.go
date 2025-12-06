package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"itab-backend/internal/database"
	"itab-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// SyncList 获取用户备份列表（远程同步接口）
func SyncList(c *gin.Context) {
	userID := c.GetUint("user_id")

	var backups []models.Backup
	if err := database.DB.Select("id, name, size, sync_count, created_at, updated_at").
		Where("user_id = ?", userID).Find(&backups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取备份列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": backups})
}

// SyncDownload 下载备份数据（远程同步接口）
func SyncDownload(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的备份ID"})
		return
	}

	userID := c.GetUint("user_id")
	username, _ := c.Get("username")
	accessKeyID := c.GetUint("access_key_id")
	accessKey, _ := c.Get("access_key")

	var backup models.Backup
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&backup).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "备份不存在"})
		return
	}

	// 记录同步记录
	record := &models.SyncRecord{
		BackupName:  backup.Name,
		TransType:   "download",
		AccessKeyID: accessKeyID,
		AccessKey:   accessKey.(string),
		UserID:      userID,
	}
	database.DB.Create(record)

	// 更新同步次数
	database.DB.Model(&backup).UpdateColumn("sync_count", backup.SyncCount+1)

	// 打印操作日志
	log.Printf("[同步] 用户 %s 使用密钥 %s 下载了备份「%s」", username, accessKey, backup.Name)

	// 解析data为对象
	var backupData interface{}
	if err := json.Unmarshal([]byte(backup.Data), &backupData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解析备份数据失败"})
		return
	}

	// 返回完整的导出格式
	c.JSON(http.StatusOK, gin.H{
		"version":            "2.1",
		"exportDate":         backup.UpdatedAt,
		"passwordsEncrypted": backup.PasswordsEncrypted,
		"data":               backupData,
	})
}

// SyncUploadRequest 上传请求
type SyncUploadRequest struct {
	Name               string      `json:"name" binding:"required"` // 备份名称
	Data               interface{} `json:"data" binding:"required"` // 备份数据
	PasswordsEncrypted bool        `json:"passwordsEncrypted"`      // 密码是否加密
}

// SyncUpload 上传备份数据（远程同步接口）
func SyncUpload(c *gin.Context) {
	userID := c.GetUint("user_id")
	username, _ := c.Get("username")
	accessKeyID := c.GetUint("access_key_id")
	accessKey, _ := c.Get("access_key")

	var req SyncUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	// 序列化data为字符串
	importData := ""
	if dataJSON, err := json.Marshal(req.Data); err == nil {
		importData = string(dataJSON)
	}

	dataSize := int64(len(importData))

	// 查找是否存在同名备份
	var existingBackup models.Backup
	err := database.DB.Where("name = ? AND user_id = ?", req.Name, userID).First(&existingBackup).Error

	if err == nil {
		// 更新现有备份
		existingBackup.Data = importData
		existingBackup.Size = dataSize
		existingBackup.SyncCount++
		existingBackup.PasswordsEncrypted = req.PasswordsEncrypted
		if err := database.DB.Save(&existingBackup).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新备份失败"})
			return
		}

		// 记录同步记录
		record := &models.SyncRecord{
			BackupName:  req.Name,
			TransType:   "upload",
			AccessKeyID: accessKeyID,
			AccessKey:   accessKey.(string),
			UserID:      userID,
		}
		database.DB.Create(record)

		// 打印操作日志
		log.Printf("[同步] 用户 %s 使用密钥 %s 更新了备份「%s」", username, accessKey, req.Name)

		c.JSON(http.StatusOK, gin.H{
			"message":   "备份更新成功",
			"backup_id": existingBackup.ID,
		})
		return
	}

	// 创建新备份
	backup := &models.Backup{
		Name:               req.Name,
		Data:               importData,
		Size:               dataSize,
		SyncCount:          1,
		PasswordsEncrypted: req.PasswordsEncrypted,
		UserID:             userID,
	}

	if err := database.DB.Create(backup).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建备份失败"})
		return
	}

	// 记录同步记录
	record := &models.SyncRecord{
		BackupName:  req.Name,
		TransType:   "upload",
		AccessKeyID: accessKeyID,
		AccessKey:   accessKey.(string),
		UserID:      userID,
	}
	database.DB.Create(record)

	// 打印操作日志
	log.Printf("[同步] 用户 %s 使用密钥 %s 创建了备份「%s」", username, accessKey, req.Name)

	c.JSON(http.StatusOK, gin.H{
		"message":   "备份创建成功",
		"backup_id": backup.ID,
	})
}
