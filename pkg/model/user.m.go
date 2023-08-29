package model

import (
	"gorm.io/gorm"
)

const (
	BidLen = 12

	HexChars = "0123456789ABCDEFabcdef"
)

type User struct {
	gorm.Model
	BID          string `gorm:"type:varchar(15);uniqueIndex;not null"` // business ID
	Nick         string `gorm:"type:varchar(80)"`
	Country      string `gorm:"type:varchar(50)"`
	FirmName     string `gorm:"type:varchar(300);not null"`
	FirmType     uint   `gorm:"not null"`
	FirmVerified int64  `gorm:"not null;default 0"`
	Email        string `gorm:"type:varchar(120)"`
	Mobile       string `gorm:"type:varchar(20)"`
	Password     string `gorm:"type:varchar(100);not null"`
	Status       int8   `gorm:"not null;default 0"` // 0 正常用户 -1 封禁 etc 待定
	Ga           string `gorm:"varchar(100)"`
	LastEVTime   int64  `gorm:"not null"`
	LastMVTime   int64  `gorm:"not null"`
	LastGVTime   int64  `gorm:"not null"`
	Fid          string `gorm:"varchar(15);default null"`
	Admin        bool   `gorm:"not null;default false"`
}

type AppIDAddress struct {
	gorm.Model
	BID   string `gorm:"type:varchar(15);not null;uniqueIndex:bid_app"`
	AppID string `gorm:"type:varchar(6);not null;uniqueIndex:bid_app"`
	Eth   string `gorm:"type:varchar(42)"`
	// Btc   string `gorm:"type:varchar(100)"`
	// Trx   string `gorm:"type:varchar(100)"`
}
