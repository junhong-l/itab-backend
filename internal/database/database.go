package database

import (
	"log"

	"itab-backend/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库
func InitDB(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	// 自动迁移
	err = DB.AutoMigrate(
		&models.User{},
		&models.AccessKey{},
		&models.Backup{},
		&models.SyncRecord{},
	)
	if err != nil {
		return err
	}

	log.Println("数据库初始化完成")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
