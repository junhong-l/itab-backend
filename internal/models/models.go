package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Username      string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password      string    `json:"password,omitempty" gorm:"size:255;not null"`
	IsAdmin       bool      `json:"is_admin" gorm:"default:false"`
	NeedChangePwd bool      `json:"need_change_pwd" gorm:"default:false"` // 首次登录需要修改密码
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// AccessKey 密钥模型
type AccessKey struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	AccessKey string     `json:"access_key" gorm:"uniqueIndex;size:64;not null"`
	SecretKey string     `json:"secret_key" gorm:"size:128;not null"`
	UserID    uint       `json:"user_id" gorm:"not null"`
	User      User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at"`                      // nil表示永久有效
	IsExpired bool       `json:"is_expired" gorm:"default:false"` // 手动过期标记
}

// Backup 备份模型
type Backup struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	Name               string    `json:"name" gorm:"uniqueIndex;size:255;not null"` // 备份名称，唯一值
	Data               string    `json:"data,omitempty" gorm:"type:text"`           // JSON数据
	Size               int64     `json:"size"`                                      // 备份大小（字节）
	SyncCount          int       `json:"sync_count" gorm:"default:0"`               // 同步次数
	PasswordsEncrypted bool      `json:"passwords_encrypted" gorm:"default:true"`   // 密码是否加密
	UserID             uint      `json:"user_id" gorm:"not null"`
	User               User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// SyncRecord 同步记录模型
type SyncRecord struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	BackupName  string    `json:"backup_name" gorm:"size:255;not null"`
	TransType   string    `json:"trans_type" gorm:"size:20;not null"` // upload/download
	AccessKeyID uint      `json:"access_key_id"`
	AccessKey   string    `json:"access_key" gorm:"size:64"` // 使用的密钥
	UserID      uint      `json:"user_id" gorm:"not null"`
	User        User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt   time.Time `json:"created_at"`
}

// 备份数据结构
type BackupData struct {
	Partitions              []Partition    `json:"partitions"`
	Folders                 []Folder       `json:"folders"`
	Shortcuts               []Shortcut     `json:"shortcuts"`
	SearchEngines           []SearchEngine `json:"searchEngines"`
	Settings                Settings       `json:"settings"`
	Passwords               []Password     `json:"passwords,omitempty"`
	CurrentEngine           int            `json:"currentEngine,omitempty"`
	CurrentPartition        int            `json:"currentPartition,omitempty"`
	CurrentPrivatePartition int            `json:"currentPrivatePartition,omitempty"`
}

// Password 密码条目
type Password struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Notes     string `json:"notes"`
	CreatedAt int64  `json:"createdAt"`
}

// Partition 工作区/分区
type Partition struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Order     int    `json:"order"`
	IsPrivate bool   `json:"isPrivate"`
}

// Folder 文件夹
type Folder struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Collapsed   bool   `json:"collapsed"`
	PartitionID *int   `json:"partitionId"`
	IsPrivate   bool   `json:"isPrivate"`
}

// Shortcut 书签
type Shortcut struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Icon        string `json:"icon"`
	FolderID    *int   `json:"folderId"`
	PartitionID *int   `json:"partitionId"`
	IsPrivate   bool   `json:"isPrivate"`
	IsPinned    bool   `json:"isPinned"`
	Order       int    `json:"order,omitempty"`
}

// SearchEngine 搜索引擎
type SearchEngine struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
	Icon string `json:"icon"`
}

// Settings 外观设置
type Settings struct {
	BgType         string `json:"bgType"`
	GradientColor1 string `json:"gradientColor1"`
	GradientColor2 string `json:"gradientColor2"`
	GradientAngle  int    `json:"gradientAngle"`
	SolidColor     string `json:"solidColor"`
	BgImage        string `json:"bgImage"`
	IconSize       int    `json:"iconSize"`
	FolderSize     int    `json:"folderSize"`
	IconGap        int    `json:"iconGap"`
	FolderGap      int    `json:"folderGap"`
	IconRadius     int    `json:"iconRadius"`
	SearchRadius   int    `json:"searchRadius"`
	BtnRadius      int    `json:"btnRadius"`
}
