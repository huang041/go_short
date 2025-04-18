package entity

import (
	"time"

	"gorm.io/gorm"
)

// URLMapping 是 URL 縮短服務的核心實體
type URLMapping struct {
	gorm.Model
	ShortURL    *string    `json:"short_url" gorm:"column:short_url;type:varchar(255);unique"`
	OriginalURL string     `json:"original_url" gorm:"column:original_url;type:varchar(255);not null"` // 修改 unique 為 not null
	Algorithm   string     `json:"algorithm" gorm:"type:varchar(50);default:'base62'"`
	Visits      int        `json:"visits" gorm:"default:0"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" gorm:"index"`
	UserID      *uint      `json:"user_id,omitempty" gorm:"index"`
}

// TableName 指定資料表名稱
func (URLMapping) TableName() string {
	return "url_mappings"
}

// IsExpired 檢查 URL 是否已過期
func (u *URLMapping) IsExpired() bool {
	if u.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*u.ExpiresAt)
}

// IncrementVisits 增加訪問計數
func (u *URLMapping) IncrementVisits() {
	u.Visits++
}
