package models

import (
  "gorm.io/gorm"
)

type UrlMapping struct {
	gorm.Model `json:"id"`
	Rename_url string `json:"rename_url"`
	Origin_url string `json:"origin_url"`
}