package handlers

import (
	"net/http"
	"strconv"

	"itab-backend/internal/database"
	"itab-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	IsAdmin  bool   `json:"is_admin"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Password string `json:"password"`
	IsAdmin  *bool  `json:"is_admin"`
}

// ListUsers 获取用户列表
func ListUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	// 隐藏密码
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 检查用户名是否已存在
	var count int64
	database.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	user := &models.User{
		Username: req.Username,
		Password: req.Password,
		IsAdmin:  req.IsAdmin,
	}

	if err := database.DB.Create(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"message": "用户创建成功", "data": user})
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	if req.Password != "" {
		user.Password = req.Password
	}
	if req.IsAdmin != nil {
		user.IsAdmin = *req.IsAdmin
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户失败"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"message": "用户更新成功", "data": user})
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 检查是否存在备份数据
	var backupCount int64
	database.DB.Model(&models.Backup{}).Where("user_id = ?", id).Count(&backupCount)
	if backupCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该用户存在备份数据，请先删除备份数据"})
		return
	}

	// 检查是否存在密钥
	var keyCount int64
	database.DB.Model(&models.AccessKey{}).Where("user_id = ?", id).Count(&keyCount)
	if keyCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该用户存在密钥，请先删除密钥"})
		return
	}

	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

// GetUser 获取单个用户
func GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"data": user})
}
