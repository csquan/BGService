package model

import "gorm.io/gorm"

type Dict struct {
	gorm.Model
	Group uint16 `gorm:"not null"`
	Value string `gorm:"type:varchar(100);not null"`
	Desc  string `gorm:"type:varchar(500);not null"`
}
