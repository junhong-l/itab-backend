package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	logFile  *os.File
	logDir   string
	keepDays int
)

// InitLogger 初始化日志系统
// dir: 日志目录
// days: 保留天数，0表示永久保留
func InitLogger(dir string, days int) error {
	logDir = dir
	keepDays = days

	// 创建日志目录
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 创建日志文件，按日期命名
	filename := filepath.Join(dir, fmt.Sprintf("%s.log", time.Now().Format("2006-01-02")))
	var err error
	logFile, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("创建日志文件失败: %v", err)
	}

	// 同时输出到控制台和文件
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime)

	// 清理旧日志
	if keepDays > 0 {
		go cleanOldLogs()
	}

	log.Printf("日志系统初始化完成，日志目录: %s，保留天数: %d", dir, days)
	return nil
}

// cleanOldLogs 清理过期日志文件
func cleanOldLogs() {
	if keepDays <= 0 {
		return
	}

	cutoff := time.Now().AddDate(0, 0, -keepDays)

	files, err := os.ReadDir(logDir)
	if err != nil {
		log.Printf("读取日志目录失败: %v", err)
		return
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".log") {
			continue
		}

		// 从文件名解析日期 (格式: 2006-01-02.log)
		name := strings.TrimSuffix(file.Name(), ".log")
		fileDate, err := time.Parse("2006-01-02", name)
		if err != nil {
			continue // 跳过无法解析的文件
		}

		if fileDate.Before(cutoff) {
			filePath := filepath.Join(logDir, file.Name())
			if err := os.Remove(filePath); err != nil {
				log.Printf("删除日志文件失败 %s: %v", file.Name(), err)
			} else {
				log.Printf("已清理过期日志: %s", file.Name())
			}
		}
	}
}

// CleanLogs 手动清理指定天数前的日志
func CleanLogs(days int) (int, error) {
	if days < 0 {
		return 0, fmt.Errorf("天数不能为负数")
	}

	cutoff := time.Now().AddDate(0, 0, -days)
	deleted := 0

	files, err := os.ReadDir(logDir)
	if err != nil {
		return 0, fmt.Errorf("读取日志目录失败: %v", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".log") {
			continue
		}

		name := strings.TrimSuffix(file.Name(), ".log")
		fileDate, err := time.Parse("2006-01-02", name)
		if err != nil {
			continue
		}

		// days=0 表示删除所有（除了今天）
		if days == 0 {
			if name == time.Now().Format("2006-01-02") {
				continue // 保留今天的日志
			}
		} else if !fileDate.Before(cutoff) {
			continue
		}

		filePath := filepath.Join(logDir, file.Name())
		if err := os.Remove(filePath); err != nil {
			log.Printf("删除日志文件失败 %s: %v", file.Name(), err)
		} else {
			deleted++
			log.Printf("已删除日志: %s", file.Name())
		}
	}

	return deleted, nil
}

// GetLogFiles 获取所有日志文件信息
func GetLogFiles() ([]LogFileInfo, error) {
	files, err := os.ReadDir(logDir)
	if err != nil {
		return nil, fmt.Errorf("读取日志目录失败: %v", err)
	}

	var logFiles []LogFileInfo
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".log") {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		logFiles = append(logFiles, LogFileInfo{
			Name:    file.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})
	}

	// 按时间倒序排列
	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].ModTime.After(logFiles[j].ModTime)
	})

	return logFiles, nil
}

// LogFileInfo 日志文件信息
type LogFileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
}

// Close 关闭日志文件
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

// RotateLog 日志轮转（每天调用，切换到新日志文件）
func RotateLog() error {
	today := time.Now().Format("2006-01-02")
	currentLog := filepath.Join(logDir, today+".log")

	// 如果已经是今天的日志文件，不需要轮转
	if logFile != nil && logFile.Name() == currentLog {
		return nil
	}

	// 关闭旧文件
	if logFile != nil {
		logFile.Close()
	}

	// 打开新文件
	var err error
	logFile, err = os.OpenFile(currentLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("创建新日志文件失败: %v", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	// 清理旧日志
	if keepDays > 0 {
		go cleanOldLogs()
	}

	return nil
}

// StartLogRotation 启动日志轮转定时任务
func StartLogRotation() {
	go func() {
		for {
			// 计算到明天0点的时间
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 1, 0, now.Location())
			duration := next.Sub(now)

			time.Sleep(duration)
			if err := RotateLog(); err != nil {
				log.Printf("日志轮转失败: %v", err)
			}
		}
	}()
}
