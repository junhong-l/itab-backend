package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"itab-backend/internal/auth"
	"itab-backend/internal/database"
	"itab-backend/internal/logger"
	"itab-backend/internal/router"
)

// getEnvOrDefault 从环境变量获取值，如果不存在则返回默认值
func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// getEnvIntOrDefault 从环境变量获取整数值
func getEnvIntOrDefault(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}

func main() {
	// 命令行参数（优先级高于环境变量）
	user := flag.String("user", "", "管理员用户名")
	pwd := flag.String("pwd", "", "管理员密码")
	port := flag.Int("port", 0, "服务端口")
	dbPath := flag.String("db", "", "数据库路径")
	logDir := flag.String("log-dir", "", "日志目录")
	logKeepDays := flag.Int("log-keep-days", -1, "日志保留天数，0表示永久保留")
	flag.Parse()

	// 环境变量作为默认值，命令行参数优先
	finalUser := *user
	if finalUser == "" {
		finalUser = getEnvOrDefault("ITAB_USER", "")
	}

	finalPwd := *pwd
	if finalPwd == "" {
		finalPwd = getEnvOrDefault("ITAB_PWD", "")
	}

	finalPort := *port
	if finalPort == 0 {
		finalPort = getEnvIntOrDefault("ITAB_PORT", 8445)
	}

	finalDbPath := *dbPath
	if finalDbPath == "" {
		finalDbPath = getEnvOrDefault("ITAB_DB", "./data/itab.db")
	}

	finalLogDir := *logDir
	if finalLogDir == "" {
		finalLogDir = getEnvOrDefault("ITAB_LOG_DIR", "./logs")
	}

	finalLogKeepDays := *logKeepDays
	if finalLogKeepDays == -1 {
		finalLogKeepDays = getEnvIntOrDefault("ITAB_LOG_KEEP_DAYS", 3)
	}

	// 初始化日志系统
	if err := logger.InitLogger(finalLogDir, finalLogKeepDays); err != nil {
		log.Fatalf("日志系统初始化失败: %v", err)
	}
	defer logger.Close()

	// 启动日志轮转
	logger.StartLogRotation()

	// 确保数据目录存在
	if err := os.MkdirAll("./data", 0755); err != nil {
		log.Fatalf("创建数据目录失败: %v", err)
	}

	// 确保静态文件目录存在
	if err := os.MkdirAll("./static", 0755); err != nil {
		log.Fatalf("创建静态文件目录失败: %v", err)
	}

	// 初始化数据库
	if err := database.InitDB(finalDbPath); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 初始化管理员用户
	if finalUser != "" && finalPwd != "" {
		// 命令行指定了用户名密码，创建或更新用户
		if err := auth.InitMasterUser(finalUser, finalPwd, false); err != nil {
			log.Fatalf("初始化管理员用户失败: %v", err)
		}
		log.Printf("管理员用户已设置: %s", finalUser)
	} else {
		// 没有指定用户名密码，检查是否需要自动创建
		if auth.NeedCreateMasterUser() {
			masterUser := "master"
			masterPwd := auth.GenerateRandomString(12)
			if err := auth.InitMasterUser(masterUser, masterPwd, true); err != nil {
				log.Fatalf("初始化管理员用户失败: %v", err)
			}
			log.Println("========================================")
			log.Println("自动生成管理员账户:")
			log.Printf("用户名: %s", masterUser)
			log.Printf("密码: %s", masterPwd)
			log.Println("请登录后立即修改密码!")
			log.Println("========================================")
		}
	}

	// 启动服务
	r := router.SetupRouter()
	addr := fmt.Sprintf(":%d", finalPort)
	log.Printf("服务启动在 http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
