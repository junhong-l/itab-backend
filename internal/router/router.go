package router

import (
	"itab-backend/internal/handlers"
	"itab-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 设置信任的代理，这里设置为不信任任何代理（直接部署时使用）
	// 如果在 nginx 等反向代理后面，改为 r.SetTrustedProxies([]string{"127.0.0.1"})
	r.SetTrustedProxies(nil)

	// CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, x-access-key, x-secret-key")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 静态文件服务
	r.Static("/static", "./static")
	r.StaticFile("/", "./static/index.html")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")

	// 公开接口
	r.POST("/api/login", handlers.Login)

	// 远程同步接口（使用AccessKey认证）
	sync := r.Group("/api/sync")
	sync.Use(middleware.AccessKeyMiddleware())
	{
		sync.GET("/list", handlers.SyncList)
		sync.GET("/download/:id", handlers.SyncDownload)
		sync.POST("/upload", handlers.SyncUpload)
	}

	// 需要登录的接口
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// 用户信息
		api.GET("/me", handlers.GetCurrentUser)
		api.POST("/change-password", handlers.ChangePassword)

		// 密钥管理
		api.GET("/keys", handlers.ListKeys)
		api.POST("/keys", handlers.CreateKey)
		api.DELETE("/keys/:id", handlers.DeleteKey)
		api.POST("/keys/:id/expire", handlers.ExpireKey)

		// 备份管理
		api.GET("/backups", handlers.ListBackups)
		api.GET("/backups/:id", handlers.GetBackup)
		api.DELETE("/backups/:id", handlers.DeleteBackup)
		api.GET("/backups/:id/download", handlers.DownloadBackup)

		// 同步记录
		api.GET("/sync-records", handlers.ListSyncRecords)
		api.POST("/sync-records/clean", handlers.CleanSyncRecords)
		api.DELETE("/sync-records/:id", handlers.DeleteSyncRecord)
		api.GET("/sync-records/stats", handlers.GetSyncStats)

		// 管理员接口
		admin := api.Group("")
		admin.Use(middleware.AdminMiddleware())
		{
			// 用户管理
			admin.GET("/users", handlers.ListUsers)
			admin.POST("/users", handlers.CreateUser)
			admin.GET("/users/:id", handlers.GetUser)
			admin.PUT("/users/:id", handlers.UpdateUser)
			admin.DELETE("/users/:id", handlers.DeleteUser)

			// 日志管理
			admin.GET("/logs", handlers.GetLogFiles)
			admin.POST("/logs/clean", handlers.CleanLogs)
		}
	}

	return r
}
