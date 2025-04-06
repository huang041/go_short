package repository

import (
	"context"
	"errors"
	"time"

	"go_short/domain/urlshortener/entity"

	"gorm.io/gorm"
)

// PostgresURLRepository 是 URLRepository 的 PostgreSQL 實現
type PostgresURLRepository struct {
	db *gorm.DB
}

// NewPostgresURLRepository 創建一個新的 PostgreSQL URL 儲存庫
func NewPostgresURLRepository(db *gorm.DB) *PostgresURLRepository {
	return &PostgresURLRepository{
		db: db,
	}
}

// FindByShortURL 根據短 URL 查找映射
func (r *PostgresURLRepository) FindByShortURL(ctx context.Context, shortURL string) (*entity.URLMapping, error) {
	var mapping entity.URLMapping
	result := r.db.WithContext(ctx).Where("short_url = ?", shortURL).First(&mapping)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 返回 nil 而不是錯誤，表示未找到記錄
		}
		return nil, result.Error
	}
	return &mapping, nil
}

// FindByOriginalURL 根據原始 URL 查找映射
func (r *PostgresURLRepository) FindByOriginalURL(ctx context.Context, originalURL string) (*entity.URLMapping, error) {
	var mapping entity.URLMapping
	result := r.db.WithContext(ctx).Where("original_url = ?", originalURL).First(&mapping)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &mapping, nil
}

// Save 保存 URL 映射
func (r *PostgresURLRepository) Save(ctx context.Context, mapping *entity.URLMapping) error {
	return r.db.WithContext(ctx).Create(mapping).Error
}

// Update 更新 URL 映射
func (r *PostgresURLRepository) Update(ctx context.Context, mapping *entity.URLMapping) error {
	return r.db.WithContext(ctx).Save(mapping).Error
}

// FindAll 獲取所有 URL 映射
func (r *PostgresURLRepository) FindAll(ctx context.Context) ([]*entity.URLMapping, error) {
	var mappings []*entity.URLMapping
	result := r.db.WithContext(ctx).Find(&mappings)
	if result.Error != nil {
		return nil, result.Error
	}
	return mappings, nil
}

// DeleteExpired 刪除所有過期的 URL 映射
func (r *PostgresURLRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at IS NOT NULL AND expires_at < ?", time.Now()).Delete(&entity.URLMapping{}).Error
}
