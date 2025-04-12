package service

import (
	"go_short/domain/identity/repository" // 依賴 Repository 介面
	// 可能需要引入 entity
)

// IdentityService 定義了用戶領域的核心業務邏輯接口 (可選，但良好實踐)
type IdentityService interface {
	// 可以在這裡定義更複雜的業務操作，例如：
	// ChangeUserPassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
	// DeactivateUser(ctx context.Context, userID uint) error
}

// identityService 是 IdentityService 的具體實現
type identityService struct {
	userRepo repository.UserRepository
	// 可能注入其他依賴，如外部服務、策略物件等
}

// NewIdentityService 創建 identityService 實例
func NewIdentityService(userRepo repository.UserRepository) IdentityService {
	return &identityService{
		userRepo: userRepo,
	}
}

// --- 在這裡實現 IdentityService 介面中定義的方法 ---
// 例如 (未來可能實現):
// func (s *identityService) ChangeUserPassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
//     user, err := s.userRepo.FindByID(ctx, userID)
//     // ... 檢查舊密碼 ...
//     // ... 檢查新密碼複雜度 ...
//     // ... 更新密碼雜湊 ...
//     // ... 保存使用者 ...
//     return nil
// }
