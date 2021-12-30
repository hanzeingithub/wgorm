package test

import (
	"time"

	"gorm.io/gorm"
)

type FileInfo struct {
	ID               uint64 `gorm:"column:id"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	FatherPath       string         `gorm:"column:father_path"`
	UserId           int64          `gorm:"column:user_id"`
	FileURL          string         `gorm:"column:file_url"`
	ImageURL         string         `gorm:"column:image_url"`
	Size             string         `gorm:"column:size"`
	CustomerFileName string         `gorm:"column:customer_file_name"`
	FileName         string         `gorm:"column:file_name"`
	IsDir            bool           `gorm:"column:is_dir"`
	IsDirLocked      bool           `gorm:"column:is_dir_locked"`
	FileType         string         `gorm:"column:file_type"`
	IsDel            bool           `gorm:"column:is_del"`
	IsProtected      bool           `gorm:"column:is_protected"`
	AccessPassword   *string        `gorm:"column:access_password"`
}

func (info *FileInfo) TableName() string {
	return "file_info"
}

type WhereFileInfo struct {
	Id          *int64  `sql_field:"id"`
	UserId      *int64  `sql_field:"user_id"`
	FatherURL   *string `sql_field:"father_path"`
	FileName    *string `sql_field:"file_name"`
	FileType    *string `sql_field:"file_type"`
	IsDel       *bool   `sql_field:"is_del"`
	IsProtected *bool   `sql_field:"is_protected"`
	FileURL     *string `sql_field:"file_url"`
	IsDir       *bool   `sql_field:"is_dir"`
}

type UpdateFileInfo struct {
	FatherURL   *string `sql_field:"father_path"`
	FileName    *string `sql_field:"customer_file_name"`
	FileType    *string `sql_field:"file_type"`
	IsDel       *bool   `sql_field:"is_del"`
	IsProtected *bool   `sql_field:"is_protected"`
	FileURL     *string `sql_field:"file_url"`
}
