package repository

import (
	"context"

	"go_short/domain/identity/entity" // 引入 User 實體
)

// UserRepository 定義了使用者資料的存取操作介面
type UserRepository interface {
	// Create 創建一個新使用者
	Create(ctx context.Context, user *entity.User) error

	// FindByID 根據 ID 查找使用者
	FindByID(ctx context.Context, id uint) (*entity.User, error)

	// FindByUsername 根據使用者名稱查找使用者
	FindByUsername(ctx context.Context, username string) (*entity.User, error)

	// FindByEmail 根據電子郵件查找使用者
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// Update 更新使用者資訊 (例如 LastLogin, IsActive)
	Update(ctx context.Context, user *entity.User) error

	// Delete 標記刪除使用者 (如果使用軟刪除) 或永久刪除
	// Delete(ctx context.Context, id uint) error
}

// 可以在此文件中添加其他 Identity 相關的 Repository 介面，
// 例如 CredentialRepository, ProfileRepository 等 (如果需要的話)
