package wgorm

import (
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}
