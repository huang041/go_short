package service

import (
	"context"
	"errors"
	"log" // 引入 log

	"go_short/domain/identity/repository" // 依賴 Repository 介面
	// 可能需要引入 entity
)

// 自訂領域錯誤
var ErrUserNotFound = errors.New("service: user not found")
var ErrServiceInternal = errors.New("service: internal error") // 通用內部錯誤

// IdentityService 定義了用戶領域的核心業務邏輯接口 (可選，但良好實踐)
type IdentityService interface {
	// ActivateUser 啟用指定 ID 的使用者帳號
	ActivateUser(ctx context.Context, userID uint) error
	// DeactivateUser 停用指定 ID 的使用者帳號
	DeactivateUser(ctx context.Context, userID uint) error

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

// --- 實作 IdentityService 介面 ---

// ActivateUser 啟用指定 ID 的使用者帳號
func (s *identityService) ActivateUser(ctx context.Context, userID uint) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("Error finding user %d for activation: %v", userID, err)
		return ErrServiceInternal // 不暴露內部錯誤細節
	}
	if user == nil {
		return ErrUserNotFound
	}

	// 如果已經是啟用狀態，可以直接返回 nil，或返回特定提示（取決於業務需求）
	if user.IsActive {
		log.Printf("User %d is already active.", userID)
		return nil // 或者 return errors.New("service: user already active")
	}

	// 更新狀態
	user.IsActive = true
	if err := s.userRepo.Update(ctx, user); err != nil {
		log.Printf("Error updating user %d status for activation: %v", userID, err)
		return ErrServiceInternal
	}

	log.Printf("User %d activated successfully.", userID)
	// (可選) 發布領域事件 UserActivatedEvent
	// ...

	return nil
}

// DeactivateUser 停用指定 ID 的使用者帳號
func (s *identityService) DeactivateUser(ctx context.Context, userID uint) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		log.Printf("Error finding user %d for deactivation: %v", userID, err)
		return ErrServiceInternal
	}
	if user == nil {
		return ErrUserNotFound
	}

	// 如果已經是停用狀態，可以直接返回
	if !user.IsActive {
		log.Printf("User %d is already inactive.", userID)
		return nil
	}

	// 執行業務規則（例如，不能停用自己或特定管理員）
	// ...

	// 更新狀態
	user.IsActive = false
	if err := s.userRepo.Update(ctx, user); err != nil {
		log.Printf("Error updating user %d status for deactivation: %v", userID, err)
		return ErrServiceInternal
	}

	log.Printf("User %d deactivated successfully.", userID)
	// (可選) 發布領域事件 UserDeactivatedEvent
	// ...

	return nil
}

// --- 在這裡實現 IdentityService 介面中定義的其他方法 (未來可能實現) ---
