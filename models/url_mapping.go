package models

import (
	"gorm.io/gorm"
)

type UrlMapping struct {
	gorm.Model `json:"id"`
	Rename_url *string `json:"rename_url" gorm:"type:varchar(255);unique"`
	Origin_url string  `json:"origin_url" gorm:"type:varchar(255);unique"`
	Algorithm  string  `json:"algorithm" gorm:"type:varchar(50);default:'base62'"`
}
